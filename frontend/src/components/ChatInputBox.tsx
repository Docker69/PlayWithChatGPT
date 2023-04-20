import { FunctionComponent, memo, useContext, useEffect, useState } from "react";
import { Box, IconButton, TextField } from "@mui/material";
import { ChatContext } from "../context/ChatProvider";
import { SEND_CHAT_PROMPT, USER_ROLE } from "../global/ChatProviderConstants";
import { sendChatPrompt } from "../api/chatAPI";

const ChatInputBox: FunctionComponent = memo(() => {
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

    //wait for sendChatPrompt to return before dispatching
    sendChatPrompt(payload).then(({ success, response }) => {
      console.info("Response from sendChatPrompt: ", { success, response });
      success && dispatch({ type: SEND_CHAT_PROMPT, payload: response });
    });

    setPrompt("");
  };

  useEffect(() => {
    console.log(
      "ChatInputBox useEffect: ",
      state.activeChatSession.id
    );
  }, [state.activeChatSession.id]);

  return (
    <Box display="flex" justifySelf="flex-start">
      <TextField
        //disable this input if active session is empty
        disabled={state.activeChatSession.id === ""}
        id="outlined-basic"
        label="Start typing ....."
        variant="outlined"
        fullWidth
        sx={{ ml: 2, mr: 1, mb: 1 }}
        value={prompt}
        onChange={(v) => setPrompt(v.target.value)}
      />
      <Box sx={{ p: 0.5, mr: 3, justifySelf: "flex-end" }}>
        <IconButton
          onClick={handleSendMessage}
          //disable this input if active session is empty
          disabled={state.activeChatSession.id === ""}
        >
          <img alt="" src="/paperairplane.svg" />
        </IconButton>
      </Box>
    </Box>
  );
});

export default ChatInputBox;
