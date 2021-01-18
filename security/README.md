# Kafka Security

## 1. [SSL Encryption](https://docs.confluent.io/platform/current/security/security_tutorial.html#generate-the-keys-and-certificates)

### 1.1 Manual

#### Generate keys and certificates
```console
$ keytool -genkey -keyalg RSA -alias selfsigned -keystore keystore.jks -storepass my-password -validity 365 -keysize 4096 -storetype pkcs12 -dname "CN=localhost"
```

- dname is X500 distinguished name that used to identify entities of certificate


To view certificate inspection
```console
$ keytool -list -v -keystore keystore.jks -storepass my-password
```

#### Create your own Certificate Authority (CA)

Private key certificate generation
```console
$ openssl req -newkey rsa:4096 -nodes -keyout ca-key -x509 -days 365 -out ca-cert -subj "/CN=Kafka-Security-CA"
```

Create a trust store for the Kafka server

```console
$ keytool -keystore keystore.truststore.jks -alias CARoot -import -file ca-cert -storepass my-password -keypass my-password -noprompt
```

Create a trust store for Kafka client
```console
$ keytool -keystore kafka.client.truststore.jks -alias CARoot -importcert -file ca-cert -storepass my-password -keypass my-password -noprompt
```

#### Sign the certificate

Create a certificate signing request
```console
$ keytool -alias selfsigned -certreq -file csr-file -keystore keystore.jks -storepass my-password -keypass my-password
```

Certification Authority signs CSR
```console
$ openssl x509 -req -CA ca-cert -CAkey ca-key -in csr-file -out cert-signed -days 365 -CAcreateserial -passin pass:my-password
```

Print cert signed file
```console
$ keytool -printcert -v -file cert-signed
```

Add public cert to Keystore
```console
$ keytool -keystore keystore.jks -alias CARoot -import -file ca-cert -storepass my-password -keypass my-password -noprompt
```

Add cert signed file to Keystore
```console
$ keytool -keystore keystore.jks -import -file cert-signed -storepass my-password -keypass my-password -noprompt
```

#### Keystore inspection

```console
$ keytool -list -v -keystore kafka.client.truststore.jks
```

#### Test SSL connection
```console
$ openssl s_client -debug -connect localhost:9093 -tls1
```


## 1.2 Bash shell script

Run
```console
$ ./secrets/start.sh
```

## 1.3 Test 
Producer
```console
$ docker-compose exec kafka-cluster kafka-console-producer --broker-list localhost:19092 --topic kafka-security --producer.config /etc/kafka/secrets/host.producer.ssl.config
```

Consumer
```console
$ docker-compose exec kafka-cluster kafka-console-consumer --broker-list localhost:19092 --topic kafka-security --consumer.config /etc/kafka/secrets/host.consumer.ssl.config

```

## 2. Authorization using ACLs

Kafka ACLs were stored in Zookeeper and must secure Zookeeper. The default behavior is that no want can access the resource except a superuser
You can write your authorizer (AD, LDAP, database, Kafka)

Clients can be identified with:
- SCRAM username `alice`
- Kerberos Principal `kafka-client@hostname.com`
- Client certificate `CN=quickstart.confluent.io,OU=TEST,O=Sales,L=PaloAlto,ST=Ca,C=US`
- LDAP authorizer(commercial) integrates with RBAC

ACLs manages cluster independently like Kafka Connect, Confluent Schema Registry, KsqlDB.
ACLs control principals on Kafka resources. 

Principals are entities that can be authenticated by the authorizer. A principal is identified on security protocols (mTLS, GSSAPI, PLAIN) It has 2 types: **users** and **group** (LDAP only).
Some examples: `User:admin`, `Group:developers`, or `User:CN=quickstart.confluent.io,OU=TEST,O=Sales,L=PaloAlto,ST=Ca,C=US`


Enable ACLs in `config/server-ssl.properties`
```yaml
authorizer.class.name=kafka.security.authorizer.AclAuthorizer
```

Define Super Users for Kafka 
```yaml
super.users=User:CN=localhost;User:CN=root
```

#### [**Kafa ACLs API usage**](https://docs.confluent.io/platform/current/kafka/authorization.html#using-acls)

Add new ACLs
```console
$ kafka-acls --bootstrap-server localhost:9092 --command-config /tmp/admin.conf \ 
 --add --allow-principal User:producer \ 
 --allow-host '*' \ 
 --operation Read --operation Write \ 
 --topic test
```

Remove ACLs
```console
$ kafka-acls --authorizer-properties zookeeper.connect=localhost:2181 \
   --remove --allow-principal User:producer \
   --producer --topic test
```  

List ACLs
```console
$ kafka-acls --bootstrap-server localhost:9092 --command-config /tmp/admin.conf --list --topic test
```

### Example SASL/SCRAM

1. Create SCRAM credentials
```console
$ kafka-configs --zookeeper localhost:2181 --alter --add-config 'SCRAM-SHA-256=[password=admin-secret],SCRAM-SHA-512=[password=admin-secret]' --entity-type users --entity-name admin
```

```console
$ kafka-configs --zookeeper localhost:2181 --alter --add-config 'SCRAM-SHA-256=[password=client-secret],SCRAM-SHA-512=[password=client-secret]' --entity-type users --entity-name client
```


2. Broker configuration enable SASL/SCRAM mechanism

```yaml
# List of enabled mechanisms, can be more than one
sasl.enabled.mechanisms=SCRAM-SHA-256

# Specify one of the SASL mechanisms
sasl.mechanism.inter.broker.protocol=SCRAM-SHA-256

# Configure SASL_SSL if SSL encryption is enabled, otherwise configure SASL_PLAINTEXT
security.inter.broker.protocol=SASL_SSL
```

3. Listeners
```yaml
# With SSL encryption
listeners=SASL_SSL://kafka1:9093
advertised.listeners=SASL_SSL://localhost:9093

# Without SSL encryption
listeners=SASL_PLAINTEXT://kafka1:9093
advertised.listeners=SASL_PLAINTEXT://localhost:9093
```

4. Configure JAAS broker
```
listener.name.sasl_ssl.scram-sha-256.sasl.jaas.config=org.apache.kafka.common.security.scram.ScramLoginModule required
   username="admin"
   password="admin-secret";
```
