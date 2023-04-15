import React, { createContext, useReducer, ReactNode } from "react";
import {
  ActionType,
  NEW_CHAT_SESSION,
  SEND_CHAT_PROMPT,
  SYSTEM_ROLE,
  StateType,
} from "./storeConstants";

const initialState: StateType = {
  chatSessions: [],
  activeChatSession: { id: "", role: "", messages: [] },
};

const store = createContext<{
  state: StateType;
  dispatch: React.Dispatch<ActionType>;
}>({
  state: initialState,
  dispatch: () => {},
});

const StateProvider = ({ children }: { children: ReactNode }) => {
  const [state, dispatch] = useReducer(
    (state: StateType, action: ActionType) => {
      switch (action.type) {
        case NEW_CHAT_SESSION: {
          console.info(`${NEW_CHAT_SESSION}: `, action.payload);
          let newSession = {
            ...action.payload,
            messages: [{ role: SYSTEM_ROLE, content: action.payload.role }],
          };
          return {
            ...state,
            chatSessions: [...state.chatSessions, newSession],
            activeChatSession: newSession,
          };
        }
        case SEND_CHAT_PROMPT: {
          console.info(`${SEND_CHAT_PROMPT}: `, action.payload);
          return { ...state, activeChatSession: action.payload };
        }
        // Handle additional action types here
        default:
          return state;
      }
    },
    initialState
  );

  return (
    <store.Provider value={{ state, dispatch }}>{children}</store.Provider>
  );
};

export { store, StateProvider };
