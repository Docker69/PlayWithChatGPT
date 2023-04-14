import React, { createContext, useReducer, ReactNode } from "react";

type ChatSession = {
  id: string;
  messages: string[];
};

type StateType = {
  chatSessions: ChatSession[];
};

type ActionType =
  | { type: "NEW_CHAT_SESSION"; payload: ChatSession }
  | { type: "NEW_CHAT_PROMPT"; payload: ChatSession }
  // Add action send prompt here
  
  ;

const initialState: StateType = {
  chatSessions: [],
};

const store = createContext<{ state: StateType, dispatch: React.Dispatch<ActionType> }>({
  state: initialState,
  dispatch: () => {},
});

const StateProvider = ({ children }: { children: ReactNode }) => {
  const [state, dispatch] = useReducer(
    (state: StateType, action: ActionType) => {
      switch (action.type) {
        case "NEW_CHAT_SESSION":
          {
            console.info('Adding new chat session: ', action.payload)
            return { ...state, chatSessions: [...state.chatSessions, action.payload] };
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
