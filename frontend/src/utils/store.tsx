import React, { createContext, useReducer, ReactNode } from "react";
import { ActionType, NEW_CHAT_SESSION, SEND_CHAT_PROMPT, StateType } from "./storeConstants";
import { sendChatPrompt } from "../api/chatAPI";

const initialState: StateType = {
  chatSessions: [],
  chatPrompt: {id: "", prompt: ""},
};

const store = createContext<{ state: StateType, dispatch: React.Dispatch<ActionType> }>({
  state: initialState,
  dispatch: () => {},
});

const StateProvider = ({ children }: { children: ReactNode }) => {
  const [state, dispatch] = useReducer(
    (state: StateType, action: ActionType) => {
      switch (action.type) {
        case NEW_CHAT_SESSION:
          {
            console.info('Adding new chat session: ', action.payload);
            return { ...state, chatSessions: [...state.chatSessions, action.payload] };
          }
          case SEND_CHAT_PROMPT:
            {
              console.info('Send prompt: ', action.payload);
              sendChatPrompt(action.payload);
              return { ...state, chatPrompt: {id: "", prompt: ""} };
            }
          // Handle additional action types here
        default:
          return state;
      }
    },
    initialState
  );

  return <store.Provider value={{ state, dispatch }}>{children}</store.Provider>;
};

export { store, StateProvider };
