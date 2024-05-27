import sys
#sys.path.append('../')
import smbus
import time
import os
import RPi.GPIO as GPIO
from connectors.connector import conn_controller, Error
import sys
import subprocess
import subprocess

# Set the DISPLAY environment variable
os.environ["DISPLAY"] = ":0"

# Disable screen blanking
subprocess.Popen("xset s off", shell=True)

# Disable DPMS (Display Power Management Signaling)
subprocess.Popen("xset -dpms", shell=True)

# Disable screen blanking
subprocess.Popen("xset s noblank", shell=True)

# Launch Chromium in kiosk mode with the specified URL
subprocess.Popen("chromium-browser --noerrdialogs --disable-infobars --kiosk --incognito https://app.bingolive.games", shell=True)

GPIO.setmode(GPIO.BOARD)

heat = 38
cooling = 37
blower = 36

def turnRelayOn(chanel):
    GPIO.output(chanel, 0)

def turnRelayOff(chanel):
    GPIO.output(chanel, 1)

GPIO.setup(heat, GPIO.OUT)
GPIO.setup(cooling, GPIO.OUT)
GPIO.setup(blower, GPIO.OUT)

turnRelayOff(heat)
turnRelayOff(cooling)
turnRelayOff(blower)

lastTemp = -1

i2c = smbus.SMBus(1)
addr = 0x44
i2c.write_byte_data(addr, 0x23, 0x34)
time.sleep(0.5)

def getTempAndHumidity():
    i2c.write_byte_data(addr, 0xe0, 0x0)
    data = i2c.read_i2c_block_data(addr, 0x0, 6)
    rawT = ((data[0]) << 8) | (data[1])
    rawR = ((data[3]) << 8) | (data[4])
    temp = -45 + rawT * 175 / 65535
    tempf = (temp * (9/5)) + 32 
    RH = 100 * rawR / 65535
    return {
        "temp": float("%.1f" % tempf),
        "humidity": float("%.1f" % RH)
    }


    
def getThermoState(db):
    conn = db.cursor(prepared=True, dictionary=True)
    query = "SELECT * FROM thermostatProperties;"
    conn.execute(query)
    return conn.fetchall()[0]

def updateThermoState(db, setTemp, indicatedTemp, heat, cooling, blower):
    conn = db.cursor(prepared=True, dictionary=True)
    query = "UPDATE thermostatProperties SET timeStamp = now(), indicatedTemp = %s, heat = %s, cooling = %s, blower = %s;"
    conn.execute(query, [indicatedTemp, heat, cooling, blower])
    
    query = "INSERT INTO thermostatState(setTemp, indicatedTemp, heat, cooling, blower) VALUES ( %s, %s, %s, %s, %s );"
    conn.execute(query, [setTemp, indicatedTemp, heat, cooling, blower])
    db.commit()
    conn.close()

def getRelayState():
    heatState = True if GPIO.input(heat) == 0 else False
    coolingState = True if GPIO.input(cooling) == 0 else False
    blowerState = True if GPIO.input(blower) == 0 else False
    return {
        "heat": heatState,
        "cooling": coolingState,
        "blower": blowerState
        
    }

def decideMotions(temp, setTemp):
    if(int(temp) >= 2 + setTemp):
        return {
            "cooling": True,
            "heat": False,
            "blower": True
        }
    elif(int(temp) <= setTemp - 2):
        return {
            "cooling": False,
            "heat": True,
            "blower": True
        }
    elif(int(temp) == setTemp - 1 or int(temp) == 1 + setTemp):
        return {
            "cooling": False,
            "heat": False,
            "blower": True
        }
    else:
        return {
            "cooling": False,
            "heat": False,
            "blower": False
        }

def doThermoMotions(motions):
    if(motions["cooling"] and motions["heat"]):
        print("Heat and Cool conflict. Emergency stop.")
        sys.exit("Heat and Cool conflict. Emergency stop.")
    if(motions["cooling"]):
        turnRelayOn(cooling)
    else:
        turnRelayOff(cooling)
    if(motions["heat"]):
        turnRelayOn(heat)
    else:
        turnRelayOff(heat)
    if(motions["blower"]):
        turnRelayOn(blower)
    else:
        turnRelayOff(blower)
        
def getPreviousTemps(db, limit):
    conn = db.cursor(dictionary=True, prepared=True)
    query = f"SELECT indicatedTemp as temp FROM thermostatState ORDER BY pID DESC LIMIT {limit};"
    conn.execute(query)
    return conn.fetchall()
        
def calculateStableTemp(currentTemp, previousTemps):
    tempTotal = currentTemp
    for temp in previousTemps:
        tempTotal += float(temp["temp"])
    return int(tempTotal / (len(previousTemps) + 1))

def getMotions(db):
    conn = db.cursor(dictionary=True, prepared=True)
    query = "SELECT cooling, heat, blower FROM thermostatProperties ORDER BY timeStamp DESC LIMIT 1;"
    conn.execute(query)
    motions = conn.fetchone()
    conn.fetchall()
    
    cleanedMotions = {}
    
    for motion in motions:
        if(motions[motion] == 1):
            cleanedMotions[motion] = True
        else:
            cleanedMotions[motion] = False
        
    return cleanedMotions
    
try:
    while True:
        db = conn_controller.get_db_conn()
        setTemp = getThermoState(db)["setTemp"]
        temps = getTempAndHumidity()
        previousTemps = getPreviousTemps(db, 5)
        stableTemp = calculateStableTemp(temps["temp"], previousTemps)
        #print("stable", stableTemp)
        #motions = decideMotions(stableTemp, setTemp)
        motions = getMotions(db)
        print(motions)
        print(temps)
        print(setTemp)
        doThermoMotions(motions)
        updateThermoState(db, setTemp, temps["temp"], motions["heat"], motions["cooling"], motions["blower"])
        conn_controller.close_conn(db)
        time.sleep(1)
        os.system("clear")
except KeyboardInterrupt:
    GPIO.cleanup()
    print("\nShuttingdown thermostat")
    
