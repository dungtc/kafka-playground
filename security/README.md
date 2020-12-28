# Kafka Security

## 1. [SSL Encryption](https://docs.confluent.io/platform/current/security/security_tutorial.html#generate-the-keys-and-certificates)

### 1.1 Manual

#### Generate keys and certificates
```
keytool -genkey -keyalg RSA -alias selfsigned -keystore keystore.jks -storepass my-password -validity 365 -keysize 4096 -storetype pkcs12 -dname "CN=localhost"
```

- dname is X500 distinguished name that used to identify entities of certificate


To view certificate inspection
```
keytool -list -v -keystore keystore.jks -storepass my-password
```

#### Create your own Certificate Authority (CA)

Private key certficiate generation
```
openssl req -newkey rsa:4096 -nodes -keyout ca-key -x509 -days 365 -out ca-cert -subj "/CN=Kafka-Security-CA"
```

Create trust store for Kafka server

```
keytool -keystore keystore.truststore.jks -alias CARoot -import -file ca-cert -storepass my-password -keypass my-password -noprompt
```

Create trust store for Kafka client
```
keytool -keystore kafka.client.truststore.jks -alias CARoot -importcert -file ca-cert -storepass my-password -keypass my-password -noprompt
```

#### Sign the certificate

Create certificate signing request
```
keytool -alias selfsigned -certreq -file csr-file -keystore keystore.jks -storepass my-password -keypass my-password
```

Certification Authority signs CSR
```
openssl x509 -req -CA ca-cert -CAkey ca-key -in csr-file -out cert-signed -days 365 -CAcreateserial -passin pass:my-password
```

Print cert signed file
```
keytool -printcert -v -file cert-signed
```

Add public cert to keystore
```
keytool -keystore keystore.jks -alias CARoot -import -file ca-cert -storepass my-password -keypass my-password -noprompt
```

Add cert signed file to keystore
```
keytool -keystore keystore.jks -import -file cert-signed -storepass my-password -keypass my-password -noprompt
```

#### Keystore inspection

```
keytool -list -v -keystore kafka.client.truststore.jks
```

#### Test SSL connection
```
openssl s_client -debug -connect localhost:9093 -tls1
```


## 1.2 Bash shell script

Run
```
./secrets/start.sh
```

## 3. Test 
Producer
```
docker-compose exec kafka-cluster kafka-console-producer --broker-list localhost:19092 --topic kafka-security --producer.config /etc/kafka/secrets/host.producer.ssl.config
```

Consumer
```
docker-compose exec kafka-cluster kafka-console-consumer --broker-list localhost:19092 --topic kafka-security --consumer.config /etc/kafka/secrets/host.consumer.ssl.config

```
