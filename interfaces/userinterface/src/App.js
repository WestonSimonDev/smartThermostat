import styles from './App.module.css';
import React from "react";
import { store } from "./store/store";
import {startWsManager} from "./socket/socket";

import {changeFan, changeMode, changeTemp} from "./apiCaller/apiCaller";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      cool: false,
      heat: false,
      blower: false,
      mode: "auto",
      fan: "auto",
      indicatedTemp: 0,
      setTemp: 0
    };
  }

  setCool = () => {
    console.log("hi");
    store.dispatch({ type: "actions/updateCool", payload: { cool: true } });
  };

  setWarm = () => {
    console.log("hi");
    store.dispatch({ type: "actions/updateHeat", payload: { heat: true } });
  };

  updateState = () => {
    const storeState = store.getState();
    this.setState({
      cool: storeState.cool,
      heat: storeState.heat,
      blower: storeState.blower,
      mode: storeState.mode,
      fan: storeState.fan,
      indicatedTemp: storeState.indicatedTemp,
      setTemp: storeState.setTemp
    });

  };

  componentDidMount() {
    store.subscribe(this.updateState);

  }

  changeTemp = (dir) =>{
    this.setState({setTemp: this.state.setTemp + dir})
    const handleResponse = (responsePayload) =>{
      console.log(responsePayload)
      if(responsePayload.error === "successful"){

      }
    }
    changeTemp(handleResponse, dir)

  }

  updateMode = (mode) => {
    this.setState({mode: mode})
    const handleResponse = (responsePayload) =>{
      console.log(responsePayload)
      if(responsePayload.error === "successful"){

      }
    }
    changeMode(handleResponse, mode)
  }

  updateFan = (fan) => {
    this.setState({fan: fan})
    const handleResponse = (responsePayload) =>{
      console.log(responsePayload)
      if(responsePayload.error === "successful"){

      }
    }
    changeFan(handleResponse, fan)
  }

  render() {
    return (
        <div className={styles.container} data-cool={this.state.cool} data-heat={this.state.heat}>
          <div className={styles.temps}>
            <div className={styles.glass}>
              <div className={styles.tempContainer}>
                <div>
                  <span className={styles.temp}>
                    <span style={{fontSize: "20px"}}>
                      Set Temp
                    </span>
                    <span>{this.state.setTemp}</span>
                  </span>
                  <span className={styles.temp}>
                    <span style={{fontSize: "20px"}}>

                    </span>
                    <span>/</span>
                  </span>
                  <span className={styles.temp}>
                    <span style={{fontSize: "20px"}}>
                      Room Temp
                    </span>
                    <span>{this.state.indicatedTemp}</span>
                  </span>

                </div>
              </div>
            </div>
          </div>
          <div className={styles.settings}>
            <div className={styles.mode + " " + styles.defaultGlass}>
              <div onClick={() => this.updateMode("off")}
                   className={this.state.mode === "off" ? styles.active : ""}>Off
              </div>
              <div onClick={() => this.updateMode("cool")}
                   className={this.state.mode === "cool" ? styles.active : ""}>Cool
              </div>
              <div onClick={() => this.updateMode("auto")}
                   className={this.state.mode === "auto" ? styles.active : ""}>Auto
              </div>
              <div onClick={() => this.updateMode("heat")}
                   className={this.state.mode === "heat" ? styles.active : ""} style={{paddingLeft: "5px"}}>Heat
              </div>
            </div>
            <div className={styles.fan + " " + styles.defaultGlass}>
              <div>Fan:</div>
              <div onClick={() => this.updateFan("auto")}
                   className={this.state.fan === "auto" ? styles.active : ""}>Auto
              </div>
              <div onClick={() => this.updateFan("on")}
                   className={this.state.fan === "on" ? styles.active : ""}>On
              </div>

            </div>
          </div>
          <div className={styles.buttons}>
            <div className={styles.defaultGlass + " " + styles.button} onClick={() => this.changeTemp(1)}>
              ↑
            </div>
            <div className={styles.defaultGlass + " " + styles.button} style={{transform: "rotate(180deg)", color: "rgba(28,61,115,0.64)"}} onClick={() => this.changeTemp(-1)}>
              ↑
            </div>
          </div>

        </div>
    );
  }
}

export default App;
