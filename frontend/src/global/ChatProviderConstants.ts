//create constants for action types
export const NEW_CHAT_SESSION = "NEW_CHAT_SESSION";
export const SEND_CHAT_PROMPT = "SEND_CHAT_PROMPT";
export const SET_DRAWER_STATE = "SET_DRAWER_STATE";
export const SET_NEW_CHAT_DIALOG_STATE = "SET_NEW_CHAT_DIALOG_STATE";

//create constants for roles
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

export type ChatStateType = {
  chatSessions: ChatSession[];
  activeChatSession: ChatSession;
  mobileDrawerOpen: boolean;
  chatDialogOpen: boolean;
};

//define action types with imported constants
export type ChatActionType =
  | { type: typeof NEW_CHAT_SESSION; payload: ChatSession }
  | { type: typeof SEND_CHAT_PROMPT; payload: ChatSession }
  | { type: typeof SET_DRAWER_STATE; payload: boolean }
  | { type: typeof SET_NEW_CHAT_DIALOG_STATE; payload: boolean };

// Add action send prompt here

