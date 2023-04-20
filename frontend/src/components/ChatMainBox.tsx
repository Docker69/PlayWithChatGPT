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
  let bothMessages: GridChatMessagesType[] = [];
  let recievedMessages: GridChatMessagesType[] = [];

  //use effect to update state when active chat session changes or new message is sent
  useEffect(() => {
    console.debug("ChatMainBox useEffect: ", bothMessages);
  }, [state.activeChatSession.messages]);

  state.activeChatSession.messages.map((message) => {
    //split message.content into array of strings by new line or carriage return
    let lines = message.content.split(/[\n\r]+/);
    //remove empty string from array
    lines = lines.filter((line) => line !== "");
    console.debug("ChatMainBox: lines count in return message:", lines.length);
    recievedMessages.push({
      side: message.role === USER_ROLE ? "right" : "left",
      avatar: message.role !== USER_ROLE ? AVATAR : "",
      messages: lines,
    });
  });

  bothMessages = [
    //  ...sampleMessages,
    ...recievedMessages,
  ];

  return (
    <Box overflow="auto" style={{ maxHeight: '100%' }} height="100vh" justifySelf="flex-start">
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
