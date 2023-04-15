import { ChatSession } from "../utils/storeConstants";

const config = {
  serverAddr: 'http://localhost:80',
  //serverAddr: "",
};

export async function initChat(payload: ChatSession) {
  //initialize return value as ChatSession type and set to payload
  let response: ChatSession = payload;
  let success: boolean = false;

  // Initialize the chat session with the backend server
  try {
    const reply = await fetch(`${config.serverAddr}/api/init`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
    });
    // Handle only Status 200 responses from the server here otherwise throw an error
    if (reply.status !== 200) {
      throw new Error(`Responded with status ${reply.status}.`);
    }
    // Parse the response body as JSON
    response = await reply.json();
    success = true;
    console.debug("Chat initiated successfully: ", response);
  } catch (error: any) {
    // Handle any errors that occur during the request here
    console.error("Exception initiating chat, ", error.message);
  }
  // Return the chat session data
  return {success, response};
}

export async function sendChatPrompt(payload: ChatSession) {
  //initialize return value as ChatSession type and set to payload
  let response: ChatSession = payload;
  let success: boolean = false;

  // Send the prompt data to the backend server
  try {
    const reply = await fetch("/api/send-completion", {
      method: "POST",
      body: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
    });
    // Handle only Status 200 responses from the server here otherwise throw an error
    if (reply.status !== 200) {
      throw new Error(`Responded with status ${reply.status}.`);
    }
    // Parse the response body as JSON
    response = await reply.json();
    success = true;
    console.debug("Reply received: ", response);
  } catch (error: any) {
    // Handle any errors that occur during the request here
    console.error("Exception sending prompt, ", error.message);
  }
  // Return the chat session data
  return {success, response};
}
