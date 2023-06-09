import { ReactNode, createContext, useMemo, useReducer } from "react";
import {
  ChatActionType,
  ChatStateType,
  NEW_CHAT_SESSION,
  SEND_CHAT_PROMPT,
  SET_DRAWER_STATE,
  SET_NEW_CHAT_DIALOG_STATE,
  SYSTEM_ROLE,
  SET_WAIT_RESPONSE_STATE,
} from "../global/ChatProviderConstants";

const initialState: ChatStateType = {
  mobileDrawerOpen: false,
  chatDialogOpen: false,
  waitingForResponse: false,
  chatSessions: [],
  activeChatSession: { id: "", role: "", messages: [] },
};

export const ChatContext = createContext<{
  state: ChatStateType;
  dispatch: React.Dispatch<ChatActionType>;
}>({
  state: initialState,
  dispatch: () => {},
});

// Create a reducer
const reducer = (state: ChatStateType, action: ChatActionType) => {
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
    case SET_DRAWER_STATE:
      return { ...state, mobileDrawerOpen: action.payload };
    case SET_NEW_CHAT_DIALOG_STATE:
      return { ...state, chatDialogOpen: action.payload };
      case SET_WAIT_RESPONSE_STATE:
        return { ...state, waitingForResponse: action.payload };
      default:
      return state;
  }
};

// Create a provider
const ChatProvider = ({ children }: { children: ReactNode }) => {
  const [state, dispatch] = useReducer(reducer, initialState);

  const value = useMemo(
    () => ({
      state: state,
      dispatch: dispatch
    }),
    [state]
  );
  return (
    <ChatContext.Provider value={value}>
      {children}
    </ChatContext.Provider>
  );
};

export default ChatProvider;
