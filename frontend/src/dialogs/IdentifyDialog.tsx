import * as React from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import { Human, SET_HUMAN, SET_INIT_SESSION_DIALOG_STATE } from '../global/ChatProviderConstants';
import { ChatContext } from '../context/ChatProvider';
import { initSession } from '../api/chatAPI';

export default function IdentifyDialog() {
  const { state, dispatch } = React.useContext(ChatContext);
  const [nickName, setNickName] = React.useState("");

  const handleNickNameSubmit = () => {
    //event.preventDefault();
    //Send request to backend to initialize chat session
    console.info("Checking nickname");
    //init payload to empty human
    const payload : Human = { Id: "", Name: "", ChatIds: [], NickName: nickName };

    //wait for initSession to return before dispatching
    initSession(payload).then(({success, response}) => {
      console.info("Response from initSession: ", {success, response});
      //TODO: handle the case where nickname not found to create new human
      success && response.Id !== "" &&
        dispatch({ type: SET_HUMAN, payload: response });
    });

    handleCancel();
    setNickName("");
  };

  const handleCancel = () => {
    dispatch({
      type: SET_INIT_SESSION_DIALOG_STATE,
      payload: false,
    });
  };

  return (
    <div>
      <Dialog open={state.initSessionOpen} onClose={handleCancel}>
        <DialogTitle>Create New Chat Session</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Can I ask with whom am I chatting?
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="Please provide a nickname"
            type="text"
            fullWidth
            variant="standard"
            onChange={(v) => setNickName(v.target.value) }
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCancel}>Cancel</Button>
          <Button onClick={handleNickNameSubmit}>Submit</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}