import React, { Fragment } from 'react';
import * as Keychain from 'react-native-keychain';
import {
  SafeAreaView,
  StyleSheet,
  View,
  Text,
  StatusBar,
  Button,
} from 'react-native';

import { Colors } from 'react-native/Libraries/NewAppScreen';

const SERVER = 'example.com';

const styles = StyleSheet.create({
  body: {
    backgroundColor: Colors.white,
  },
  sectionTitle: {
    color: Colors.black,
    fontSize: 24,
    fontWeight: '600',
    marginTop: 32,
    paddingHorizontal: 24,
  },
});

const App = () => {
  const getCredentials = async (): Promise<Keychain.UserCredentials> => {
    const tmp = await Keychain.getInternetCredentials(SERVER);
    return tmp;
  };

  const setCredentials = async (username: string, password: string) => {
    await Keychain.setInternetCredentials(SERVER, username, password, {
      accessControl: Keychain.ACCESS_CONTROL.USER_PRESENCE,
      accessible: Keychain.ACCESSIBLE.WHEN_UNLOCKED_THIS_DEVICE_ONLY,
    });
  };

  const handleCheckCredentials = async () => {
    const credentials = await getCredentials();
    if (credentials) {
      console.log(credentials);
      try {
        await Keychain.resetInternetCredentials(SERVER);
      } catch (e) {
        console.log(e);
      }
    } else {
      console.log('No credentials stored');
      try {
        await setCredentials('foo', 'bar');
      } catch (e) {
        console.log(e);
      }
    }
  };

  return (
    <Fragment>
      <StatusBar barStyle="dark-content" />
      <SafeAreaView>
        <View style={styles.body}>
          <Text style={styles.sectionTitle}>This Space For Rent</Text>
        </View>
        <View style={styles.body}>
          <Button onPress={handleCheckCredentials} title="Check Credentials" />
        </View>
      </SafeAreaView>
    </Fragment>
  );
};

export default App;
