import { configureStore } from "@reduxjs/toolkit"

const initialState = {
    heat: false,
    cool: false,
    blower: false,
    mode: "auto",
    fan: "auto",
    indicatedTemp: 0,
    setTemp: 0
}

function appReducer (state = initialState, action) {
    const splitType = action.type.split("/")
    console.log(splitType)
    switch(splitType[0]){
        case "actions":
            switch(splitType[1]){
                case "updateCool":
                    var heat = state.heat
                    if(action.payload.cool){
                        heat = false
                    }
                    return {
                        ...state,
                        cool: action.payload.cool,
                        heat: heat

                    }
                case "updateHeat":
                    var cool = state.cool
                    if(action.payload.heat){
                        cool = false
                    }
                    return {
                        ...state,
                        cool: cool,
                        heat: action.payload.heat

                    }
                case "updateMode":
                    return {
                        ...state,
                        mode: action.payload.mode

                    }
                case "updateFan":
                    return {
                        ...state,
                        fan: action.payload.fan

                    }
                case "syncState":
                    return {
                        ...state,
                        cool: action.payload.cool,
                        heat: action.payload.heat,
                        blower: action.payload.blower,
                        indicatedTemp: action.payload.indicatedTemp,
                        setTemp: action.payload.setTemp,
                        mode: action.payload.mode,
                        fan: action.payload.fan
                    }
                default:
                    return state
            }
        default:
            return state
    }

}

export const store = configureStore({reducer: appReducer})

store.dispatch({type: "firstSlice/secondSlice"})