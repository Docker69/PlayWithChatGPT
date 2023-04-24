//create constants for action types
export const NEW_CHAT_SESSION = "NEW_CHAT_SESSION";
export const SEND_CHAT_PROMPT = "SEND_CHAT_PROMPT";
export const SET_DRAWER_STATE = "SET_DRAWER_STATE";
export const SET_NEW_CHAT_DIALOG_STATE = "SET_NEW_CHAT_DIALOG_STATE";
export const SET_INIT_SESSION_DIALOG_STATE = "SET_INIT_SESSION_DIALOG_STATE";
export const SET_WAIT_RESPONSE_STATE = "SET_WAIT_RESPONSE_STATE";
export const SET_HUMAN = "SET_HUMAN";

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
  humanId: string;
  messages: ChatMessages[];
};

export type Human = {
  id: string;
  name: string;
  nickName: string;
  chatIds: { id: string; role: string }[];
};

export type ChatStateType = {
  chatSessions: ChatSession[];
  activeChatSession: ChatSession;
  mobileDrawerOpen: boolean;
  chatDialogOpen: boolean;
  initSessionOpen: boolean;
  waitingForResponse: boolean;
  human: Human;
};

//define action types with imported constants
export type ChatActionType =
  | { type: typeof NEW_CHAT_SESSION; payload: ChatSession }
  | { type: typeof SEND_CHAT_PROMPT; payload: ChatSession }
  | { type: typeof SET_DRAWER_STATE; payload: boolean }
  | { type: typeof SET_NEW_CHAT_DIALOG_STATE; payload: boolean }
  | { type: typeof SET_WAIT_RESPONSE_STATE; payload: boolean }
  | { type: typeof SET_INIT_SESSION_DIALOG_STATE; payload: boolean }
  | { type: typeof SET_HUMAN; payload: Human };

// Add action send prompt here
