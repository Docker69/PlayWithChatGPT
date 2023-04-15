//create constants for action types
export const NEW_CHAT_SESSION = "NEW_CHAT_SESSION";
export const SEND_CHAT_PROMPT = "SEND_CHAT_PROMPT";

export const SYSTEM_ROLE = "system";
export const ASSISTAN_ROLE = "assistant";
export const USER_ROLE = "user";

export type ChatMessages = {
  role: string;
  content: string;
};

export type ChatSession = {
  id: string;
  role: string;
  messages: ChatMessages[];
};

export type StateType = {
  chatSessions: ChatSession[];
  activeChatSession: ChatSession;
};

//define action types with imported constants
export type ActionType =
  | { type: typeof NEW_CHAT_SESSION; payload: ChatSession }
  | { type: typeof SEND_CHAT_PROMPT; payload: ChatSession };
// Add action send prompt here
