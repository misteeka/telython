import React, {useState} from 'react';
import { StyleSheet, Button, View, SafeAreaView, Text, TextInput } from 'react-native';
import { AuthForm } from './src/AuthForm'

const App: React.FC = () => {
  
  const auth = (success: boolean) => {
    if(success) {
      console.log('SUCCESS')
    } else {
      console.log('ERROR')
    }
  }

  return (
  <SafeAreaView style={styles.container}>
    <View style={styles.navbar}>
      <Text style={styles.title}>
        checkPassword
      </Text>
    </View>
    <AuthForm onSubmit={auth}/>
  </SafeAreaView>
  )
  }

const styles = StyleSheet.create({
  navbar: {
    justifyContent: 'center',
    alignItems: 'center',
    height: 80,
    backgroundColor: 'green'
  },
  container: {
    alignContent: 'center',
    justifyContent: 'center',
  },
  title: {
    fontSize: 18,
    fontFamily: 'Helvetica, Arial, sans-serif'
  },
});

export default App;
