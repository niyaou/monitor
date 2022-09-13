import React, { useReducer } from "react";
import getDefaultConfig from './config-default-data'

const initialState = getDefaultConfig();
const context = React.createContext<any>(initialState);

const sessionState = { flag: 0x0000 }
export enum sessionFlag{
    FLAG_START_SESSION = 0X0001,
    FLAG_IDLE_SESSION = 0X0000,

}
export enum actionTypes  {
    FC_SESSION="FC_SESSION",
    FC_CONFIG="FC_CONFIG"
}
function reducer(state, action: any): any {
    switch (action.type) {
        case actionTypes.FC_SESSION:
            return {...state, ...action.payload }
        case actionTypes.FC_CONFIG:
            return { ...state, ...action.payload }
        default:
            return state
    }
}


const ContextProvider = props => {
    const [state, dispatch] = useReducer(reducer, initialState);


    return (
        <context.Provider value={{ state, dispatch }}>
            {props.children}
        </context.Provider>
    );
};

export { reducer as configDataReducer, context as configDataContext, ContextProvider as ConfigDataContextProvider };

