//create constants for action types
export const NEW_CHAT_SESSION = "NEW_CHAT_SESSION";
export const SEND_CHAT_PROMPT = "SEND_CHAT_PROMPT";

export type ChatSession = {
  id: string;
  role: string;
  messages: string[];
};

export type ChatPrompt = {
  id: string;
  prompt: string;
};

export type StateType = {
  chatSessions: ChatSession[];
  chatPrompt: ChatPrompt;
};

//define action types with imported constants
export type ActionType =
  | { type: typeof NEW_CHAT_SESSION; payload: ChatSession }
  | { type: typeof SEND_CHAT_PROMPT; payload: ChatPrompt };
// Add action send prompt here
