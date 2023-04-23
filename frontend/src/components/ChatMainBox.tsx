"use client"


import { FunctionComponent, useContext, useEffect, useRef } from "react";
import { Box } from "@mui/material";
import ChatMessages from "./ChatMessagesGrid";
import { ChatContext } from "../context/ChatProvider";
import { USER_ROLE } from "../global/ChatProviderConstants";
import { CHAT_AVATAR } from "../global/GlobalContants";
import { log } from "console";

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
  const BottomRef = useRef<HTMLDivElement>(null);
  let bothMessages: GridChatMessagesType[] = [];
  let recievedMessages: GridChatMessagesType[] = [];

  console.debug("ChatMainBox render");

  /*
  //use effect to update state when active chat session changes or new message is sent
  useEffect(() => {
    console.debug("ChatMainBox useEffect: ", state.activeChatSession.messages);
  }, [state.activeChatSession.messages]);

  //use effect to update grid when state.waitingForResponse changes
  useEffect(() => {
    //get the last StyledGrid item and add another StyledGrid item with a loading spinner
    console.debug("ChatMainBox useEffect: waitingForResponse: ", state.waitingForResponse);

  }, [state.waitingForResponse]);
  */

  useEffect(() => {
    // Scroll to the bottom of the container after rendering.
      BottomRef.current?.scrollIntoView({behavior: 'smooth'});
  }, [bothMessages]);
  
  state.activeChatSession.messages.map((message) => {
    //split message.content into array of strings by new line or carriage return
    let lines = message.content.split(/[\n\r]+/);
    //remove empty string from array
    lines = lines.filter((line) => line !== "");
    console.debug("ChatMainBox: lines count in return message:", lines.length);
    recievedMessages.push({
      side: message.role === USER_ROLE ? "right" : "left",
      avatar: message.role !== USER_ROLE ? CHAT_AVATAR : "",
      messages: lines,
    });

    return recievedMessages;
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
            last={idx === bothMessages.length - 1}
          />
        ))}
      </div>
      <div ref={BottomRef} />
    </Box>
  );
};

export default ChatMainBox;
