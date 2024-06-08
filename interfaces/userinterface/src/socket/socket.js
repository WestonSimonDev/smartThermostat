import { store } from "../store/store";
import {v4 as uuidv4} from 'uuid';


function msgHandeler(payload){
    try {
        var payload = JSON.parse(payload)
    } catch (e) {
        console.log(e)
        return
    }

    var slices = payload.action.split("/")
    if(payload.action === "response/heartBeat"){
        if(payload.payload.change === true){
            const newState = payload.payload.newState
            store.dispatch({type: "actions/syncState", payload: {
                    cool: newState.cool,
                    heat: newState.heat,
                    blower: newState.blower,
                    setTemp: newState.setTemp,
                    indicatedTemp: newState.indicatedTemp,
                    mode: newState.mode,
                    fan: newState.fan
            }})
        }

        return
    }

    switch(slices[0]){
        case "response":
            var responseId = slices[1]
            if(pendingResponses.hasOwnProperty(responseId)){
                pendingResponses[responseId](payload["payload"])
            }
            break;
        case "companies":
            switch(slices[1]){
                case "syncSelected":
                    store.dispatch({type: "companies/selectCompany", payload: payload["payload"]["selectedCompany"]})
            }
    }

}

var socket = undefined
var end = false

var sockets = {}

async function websocketManager(){
    end = false
    socket = new WebSocket(`wss://thermostat.la2.http.code.westonsimon.com/api/local/ws`)
    var socketID = uuidv4()

    socket.addEventListener("open", (event) => {
        console.log("ws opened")
        sockets[socketID] = "open"
    });

    // Listen for messages
    socket.addEventListener("message", (event) => {
        msgHandeler(event.data)
        //console.log(event.data)

    });

    socket.addEventListener("close", async (event) => {
        end = true
        sockets[socketID] = "closed"
        await new Promise(r => setTimeout(r, 3000));
        websocketManager();
        console.log("socket closed");

    })
    let heartBeat = JSON.stringify({"action": "heartBeat", "response": "heartBeatResponse"})
    var errorCount = 0

    while(true){
        if(end){
            console.log("killing heart beat")
        }
        if(sockets[socketID] === "closed"){
            break
        }
        if(end === false){
            if(socket.readyState === WebSocket.OPEN){
                socket.send(heartBeat)
                errorCount = 0
            }
        }else{
            break
        }
        await new Promise(r => setTimeout(r, 1000));
    }
}

const waitForSocketOpen = (socket, timeout = 5000) => {
    return new Promise((resolve, reject) => {
        // If the socket is already open, resolve immediately.
        if (socket.readyState === WebSocket.OPEN) {
            resolve();
        } else {
            // Attach event listeners to resolve the promise when the socket opens.
            const onOpen = () => {
                clearTimeout(timeoutId);
                socket.removeEventListener('open', onOpen);
                resolve();
            };

            // Set up a timeout to reject the promise if the socket doesn't open in time.
            const timeoutId = setTimeout(() => {
                socket.removeEventListener('open', onOpen);
                reject(new Error('WebSocket connection timeout'));
            }, timeout);

            socket.addEventListener('open', onOpen);
        }
    });
};

var pendingResponses = {}

// Updated sendWsMessage function to wait for the WebSocket to be open
export const sendWsMessage = async (payload, responseId, response) => {
    pendingResponses[responseId] = response
    try {
        // Wait for the socket to open with a timeout (e.g., 5 seconds)
        await waitForSocketOpen(socket);
        // Once the socket is open, send the message.
        socket.send(payload);
    } catch (error) {
        console.error("Failed to send message:", error.message);
        // Handle the failure (e.g., retry later, show an error message, etc.)
    }
};

export const startWsManager = (userToken) => {
    websocketManager(userToken)
}

startWsManager()
