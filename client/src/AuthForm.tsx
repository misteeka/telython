import React, { useState, Component } from 'react'
import { Text, View, TextInput, Button, StyleSheet } from 'react-native'
import Auth from './client';
import { HttpError } from './http/error';

type AuthFormData = {
    username: string,
    password: string,
}

type AuthFormProps = {
    onSubmit: (success: boolean) => void
}

export const AuthForm: React.FC<AuthFormProps> = ({onSubmit}) => {
    const [login, setLogin] = useState('')
    const [password, setPassword] = useState('')
    const [result, setResult] = useState('')

    const submitPressHandler = () => {
        Auth.CheckPassword(login, password).then((error: HttpError) => {
            if(error.error !== undefined) {
                setResult(error.error + " | " + error.message)
                onSubmit(false)
            } else {
                setResult('SUCCESS')
                onSubmit(true)
            }
        })
    }

    return (
        <View>
            <TextInput
                style={styles.input}
                onChangeText={(text) => setLogin(text)}
                placeholder='Введите имя пользователя...'
            />
            <TextInput
                style={styles.input}
                onChangeText={(text) => setPassword(text)}
                placeholder='Введите пароль...'
            />
            <Button
                title="Press me"
                onPress={submitPressHandler}
            />
            <Text>{result}</Text>
        </View>
    )
}

const styles = StyleSheet.create({
    input: {
        borderStyle: 'solid',
        borderWidth: 2,
        borderColor: 'black'
    },
    submit: {

    }
})