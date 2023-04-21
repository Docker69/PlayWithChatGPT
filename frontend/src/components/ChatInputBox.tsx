import {
  FunctionComponent,
  memo,
  useContext,
  useEffect,
  useState,
} from "react";
import { Box, IconButton, TextField } from "@mui/material";
import { ChatContext } from "../context/ChatProvider";
import {
  SEND_CHAT_PROMPT,
  SET_WAIT_RESPONSE_STATE,
  USER_ROLE,
} from "../global/ChatProviderConstants";
import { sendChatPrompt } from "../api/chatAPI";

const ChatInputBox: FunctionComponent = memo(() => {

  console.debug("ChatInputBox render");
  
  const [prompt, setPrompt] = useState("");
  const { state, dispatch } = useContext(ChatContext);

  const handleSendMessage = () => {
    console.info("Requesting to send prompt to backend");

    const updatedMessages = [
      ...state.activeChatSession.messages,
      { role: USER_ROLE, content: prompt },
    ];
    const payload = { ...state.activeChatSession, messages: updatedMessages };
    dispatch({ type: SEND_CHAT_PROMPT, payload: payload });
    dispatch({ type: SET_WAIT_RESPONSE_STATE, payload: true });

    //wait for sendChatPrompt to return before dispatching
    sendChatPrompt(payload).then(({ success, response }) => {
      console.info("Response from sendChatPrompt: ", { success, response });
      success && dispatch({ type: SEND_CHAT_PROMPT, payload: response });
      dispatch({ type: SET_WAIT_RESPONSE_STATE, payload: false });
    });

    setPrompt("");
  };

  useEffect(() => {
    console.debug(
      "ChatInputBox useEffect: id: %s, waitingForResponse: %s",
      state.activeChatSession.id,
      state.waitingForResponse
    );
  }, [state.activeChatSession.id, state.waitingForResponse]);

  return (
    <Box display="flex" justifySelf="flex-start">
      <TextField
        //conditionally add animation if waiting for response
        sx={{
          ml: 2,
          mr: 1,
          mb: 1,
        }}
        //disable this input if active session is empty
        disabled={state.activeChatSession.id === "" || state.waitingForResponse}
        id="outlined-basic"
        label={
          state.waitingForResponse ? "Thinking ....." : "Start typing ....."
        }
        variant="outlined"
        fullWidth
        value={prompt}
        onChange={(v) => setPrompt(v.target.value)}
      />
      <Box sx={{ p: 0.5, mr: 3, justifySelf: "flex-end" }}>
        <IconButton
          onClick={handleSendMessage}
          //disable this input if active session is empty
          disabled={
            state.activeChatSession.id === "" || state.waitingForResponse
          }
        >
          <img alt="" src="/paperairplane.svg" />
        </IconButton>
      </Box>
    </Box>
  );
});

export default ChatInputBox;
