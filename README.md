
# ping-kafka-mysql

## environment
  - KAFKA_HOST=test-kafka
  - KAFKA_PORT=9092
  - KAFKA_TOPIC=fruit
  - MYSQL_HOST=test-mysql
  - MYSQL_PORT=3306
  - MYSQL_USER_NAME=root
  - MYSQL_PASSWORD=1234

## docker-compose
```yml
# sample:ping-kafka-mysql
# only container(host is not) can access kakfa
services:  
  kafka-server:
    container_name: test-kafka
    environment:
      JMX_PORT: 9097
      KAFKA_ADVERTISED_HOST_NAME: test-kafka
      KAFKA_ADVERTISED_PORT: 9092
      KAFKA_BROKER_ID: 1
      KAFKA_DELETE_TOPIC_ENABLE: "true"
      KAFKA_HEAP_OPTS: -Xmx1G
      KAFKA_JMX_OPTS: -Dcom.sun.management.jmxremote=true -Dcom.sun.management.jmxremote.authenticate=false  -Dcom.sun.management.jmxremote.ssl=false
        -Dcom.sun.management.jmxremote.authenticate=false
        -Djava.rmi.server.hostname=test-kafka
      KAFKA_JVM_PERFORMANCE_OPTS: -XX:+UseG1GC -XX:MaxGCPauseMillis=20 -XX:InitiatingHeapOccupancyPercent=35
        -XX:+DisableExplicitGC -Djava.awt.headless=true
      KAFKA_LOG_CLEANER_ENABLE: "true"
      KAFKA_LOG_CLEANUP_POLICY: delete
      KAFKA_LOG_DIRS: /logs/kafka-logs-24bf1bde016a
      KAFKA_LOG_RETENTION_HOURS: 120
      KAFKA_ZOOKEEPER_CONNECT: test-zookeeper:2181
      KAFKA_ZOOKEEPER_CONNECTion_timeout_ms: 60000
    # extra_hosts:
    # - test-kafka:10.202.101.43
    image: pangpanglabs/kafka
    ports:
    - 9092
    - 9097
  zookeeper-server:
    container_name: test-zookeeper
    environment:
      ZOO_MY_ID: 1
      ZOO_SERVERS: server.1=test-zookeeper:2888:3888
    # extra_hosts:
    # - test-kafka:10.202.101.43
    image: zookeeper:3.4.9
    ports:
    - 2181
    - 2888
    - 3888
  mysql-server:
   container_name: test-mysql
   environment:
   - MYSQL_ROOT_PASSWORD=1234
   image: mysql:5.7.22
   ports:
   - 3306
   volumes:
   - ./database/mysql/:/docker-entrypoint-initdb.d
  ping-kafka-mysql-server:
    container_name: test-ping-kafka-mysql
    command: sh -c 'echo "wait kafka..." && /go/bin/wait-for.sh test-kafka:9092 test-mysql:3306 -t 36000 -- ./ping-kafka-mysql'
    depends_on:
    - kafka-server
    - mysql-server
    environment:
    - KAFKA_HOST=test-kafka
    - KAFKA_PORT=9092
    - KAFKA_TOPIC=fruit
    - MYSQL_HOST=test-mysql
    - MYSQL_PORT=3306
    - MYSQL_USER_NAME=root
    - MYSQL_PASSWORD=1234
    image: relaxed/ping-kafka-mysql
    volumes:
      - ./wait-for.sh:/go/bin/wait-for.sh
    ports:
    - 8080
version: "3"
 
```

[wait-for.sh](https://github.com/Eficode/wait-for)
