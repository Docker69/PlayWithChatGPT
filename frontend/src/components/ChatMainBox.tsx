import { FunctionComponent, useContext, useEffect } from "react";
import { Box } from "@mui/material";
import ChatMessages from "./ChatMessagesGrid";
import { ChatContext } from "../context/ChatProvider";
import { USER_ROLE } from "../global/ChatProviderConstants";

const AVATAR = "/logo50.svg";

type GridChatMessagesType = {
  side: "left" | "right";
  avatar: string;
  messages: string[];
};

/*
const sampleMessages: GridChatMessagesType[] = [
  {
    side: "left",
    avatar: AVATAR,
    messages: [
      "Hi Jenny, How r u today?",
      "Did you train yesterday",
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Volutpat lacus laoreet non curabitur gravida.",
    ],
  },
  {
    side: "right",
    avatar: "",
    messages: [
      "Great! What's about you?",
      "Of course I did. Speaking of which check this out",
    ],
  },
  {
    side: "left",
    avatar: AVATAR,
    messages: ["Im good.", "See u later."],
  },
];
*/
const ChatMainBox: FunctionComponent = () => {
  const { state } = useContext(ChatContext);
  //use effect to update state when active chat session changes or new message is sent
  useEffect(() => {
    console.log(
      "ChatMessagesGrid useEffect: ",
      state.activeChatSession.messages
    );
  }, [state.activeChatSession.messages]);

  let recievedMessages: GridChatMessagesType[] = [];
  state.activeChatSession.messages.map((message) => (
    recievedMessages.push({
      side: message.role === USER_ROLE ? "right" : "left",
      avatar: AVATAR,
      messages: [message.content],
    })
  ));

  let bothMessages: GridChatMessagesType[] = [
  //  ...sampleMessages,
    ...recievedMessages,
  ];
  return (
    <Box height="100vh" justifySelf="flex-start">
      <div>
        {bothMessages.map((message, idx) => (
          <ChatMessages
            key={idx}
            side={message.side}
            avatar={message.avatar}
            messages={[...message.messages]}
          />
        ))}
      </div>
    </Box>
  );
};

export default ChatMainBox;
