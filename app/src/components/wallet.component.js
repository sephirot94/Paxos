import React, { useEffect, useState } from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Link from '@material-ui/core/Link';
import Grid from '@material-ui/core/Grid';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import { useHistory } from 'react-router-dom';
import Alert from '@material-ui/lab/Alert';
import { InputLabel } from '@material-ui/core';

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://material-ui.com/">
        PAXOS
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    height: 200,
    width: 200
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
    backgroundColor: '#118ECB'
  },
  links: {
      color: '#118ECB',
  },
  input: {
    width: '80%'
  },
  item: {
    marginTop: 10,
    fontSize: 30
  }
}));

export default function Wallet() {
  const classes = useStyles();
  let history = useHistory();
  
  const [montoPagar, setMontoPagar] = useState(0);
  const [montoCobrar, setMontoCobrar] = useState(0);
  const [saldo, setSaldo] = useState();
  const [error, setError] = useState();

  function pagarRequest(){
    const data = {
      type: 'debit',
      ammount: montoPagar
    }
    fetch("http://localhost:8080/transactions", {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
      .then(res => res.json())
      .then(res => {
        
      })
  }

  function cobrarRequest(){
    const data = {
      type: 'debit',
      ammount: montoPagar
    }
    fetch("http://localhost:5000/api/login", {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
      .then(res => res.json())
  }

  useEffect(()=>{
    // if(loginData.isSuccess){
    //   history.push("/profile", {userData:loginData.data});
    // }
  });

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <div className={classes.paper}>
        <Avatar variant="square" src="/assets/img/wallet.png" className={classes.avatar} />
        <Grid container>
          <Grid item md={9}>
            <TextField
              variant="outlined"
              margin="normal"
              fullWidth
              id="debit"
              label="Monto"
              name="pago"
              autoFocus
              onChange={(event) => {
                const { value } = event.target;
                setMontoPagar(value);
              }}
              className={classes.input}
            />
          </Grid>
          <Grid item md={3}>
            <Button
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
              onClick={pagarRequest}
            >
              Pagar
            </Button>
          </Grid>
          <Grid item md={9}>
            <TextField
              variant="outlined"
              margin="normal"
              fullWidth
              id="credit"
              label="Monto"
              name="credit"
              autoFocus
              onChange={(event) => {
                const { value } = event.target;
                setMontoCobrar(value);
              }}
              className={classes.input}
            />
          </Grid>
          <Grid item md={3}>
            <Button
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
              onClick={cobrarRequest}
            >
              Cobrar
            </Button>
          </Grid>
          <Grid item md={3} >
            <InputLabel
                margin="normal"
                id="saldoLbl"
                name="saldoLbl"
                className={classes.item}
                margin= 'dense'
              >
                Saldo:
                </InputLabel>
          </Grid>
          <Grid item md={9}>
            <TextField
                variant="outlined"
                margin="normal"
                fullWidth
                id="saldo"
                label={saldo}
                name="saldo"
                disabled
                className={classes.input}
              />
          </Grid>
        </Grid>
        {error ? <Alert severity="error">error</Alert> : null}
      </div>
      <Box mt={8}>
        <Copyright />
      </Box>
    </Container>
  );
}