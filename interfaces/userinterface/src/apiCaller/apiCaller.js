import {sendWsMessage} from "../socket/socket";
import {v4 as uuidv4} from 'uuid';

export const changeTemp = async (callback, dir) => {
    var responseId = uuidv4()
    const payload = {"action": "temps/change", "payload": {"direction": dir}, "response": responseId}
    sendWsMessage(JSON.stringify(payload),  responseId, callback)
}

export const changeMode = async (callback, mode) => {
    var responseId = uuidv4()
    const payload = {"action": "mode/change", "payload": {"mode": mode}, "response": responseId}
    sendWsMessage(JSON.stringify(payload),  responseId, callback)
}

export const changeFan = async (callback, mode) => {
    var responseId = uuidv4()
    const payload = {"action": "fan/change", "payload": {"fan": mode}, "response": responseId}
    sendWsMessage(JSON.stringify(payload),  responseId, callback)
}