import { ChatPrompt, ChatSession } from "../utils/storeConstants";

const config = {
    //serverAddr: 'http://localhost:3000',
    serverAddr: '',
  };
  
export async function initChat(payload: ChatSession) {
    //initialize return value as ChatSession type and set to payload
    let chatSession: ChatSession = payload;

    // Initialize the chat session with the backend server
    try {
        const response = await fetch(`${config.serverAddr}/api/init`, {
            method: 'POST',
            body: JSON.stringify(payload),
            headers: {
                'Content-Type': 'application/json'
            }
        });
        // Handle only Status 200 responses from the server here otherwise throw an error
        if (response.status !== 200) {
            throw new Error(`Responsed with status ${response.status}.`);
        }
        // Parse the response body as JSON
        chatSession = await response.json();
        console.debug('Chat initiated successfully: ', chatSession);

    } catch (error: any) {
        // Handle any errors that occur during the request here
        console.error('Exception initiating chat, ', error.message);
    }
    // Return the chat session data
    return chatSession;
  }

export async function sendChatPrompt(payload: ChatPrompt) {
    // Send the prompt data to the backend server
    try {
        const response = await fetch('/api/send-prompt', {
            method: 'POST',
            body: JSON.stringify(payload),
            headers: {
                'Content-Type': 'application/json'
            }
        });
        // Handle the response from the server here if necessary
        console.debug('Prompt sent successfully: ', response);
    } catch (error) {
        // Handle any errors that occur during the request here
        console.error('Error sending prompt: ', error);
    }
  }