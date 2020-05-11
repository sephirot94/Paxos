import React, {useEffect} from 'react';
import logo from './logo.svg';
import './App.css';
import axios from 'axios'
import {Card, Descriptions,} from 'antd'

export default function App() {
    const [transactionHistory, setTransactionHistory] = useState([{"id_transaction":"5577006791947779410","type":"credit","ammount":10000,"date":"2020-04-26 21:24:08.609019 -0300 -03 m=+3.010234986"},{"id_transaction":"5577006791947779410","type":"credit","ammount":10000,"date":"2020-04-26 21:24:08.609019 -0300 -03 m=+3.010234986"},{"id_transaction":"5577006791947779410","type":"credit","ammount":10000,"date":"2020-04-26 21:24:08.609019 -0300 -03 m=+3.010234986"},{"id_transaction":"5577006791947779410","type":"credit","ammount":10000,"date":"2020-04-26 21:24:08.609019 -0300 -03 m=+3.010234986"}])
    const [transaction, setTransaction] = useState({"id_transaction":"5577006791947779410","type":"credit","ammount":10000,"date":"2020-04-26 21:24:08.609019 -0300 -03 m=+3.010234986"})
    const [accountBalance, setAccountBalance] = useState(500)
    // const [modalOpened, setModalOpened]= useState(false)
    const api = axios.create({
        baseURL: 'http://localhost:8080/',
        headers : {
            'Content-Type' : 'application/json'
        }
    });

    const getTransactionHistory = () => {
        return api.get('/transactions')
    }

    const getAccountBalance = () => {
        return api.get('/account')
    }

    const getTransaction = (id) => {
        return api.get(`/miami/accounts/${id}`)
    }

    const generateTransaction = (type, ammount) => {
        const body = {
            type: type,
            ammount: ammount
        };
        return api.post('/transactions', body)
    }

    // FE functionality
    // useEffect(() => {
    //     getTransactionHistory().then((response) => {
    //         setTransactionHistory(response)
    //     });
    // });

    // useEffect(() => {
    //     getAccountBalance().then((response) => {
    //         setAccountBalance(response);
    //     });
    // });


    // Render
    return (
        <div>
            <div>
                <h1>
                    Ledger
                </h1>
            </div>
            <p>Current Balance: {accountBalance}</p>
            { transactionHistory.forEach((transaction) => {
                return (
                    <Card>
                        <Descriptions title="Transaction">
                            <Descriptions.Item label="Type">{transaction.type}</Descriptions.Item>
                            <Descriptions.Item label="Ammount">{transaction.ammount}</Descriptions.Item>
                        </Descriptions>
                    </Card>
                    )
            })}
        </div>

    )
}
