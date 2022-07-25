FROM openjdk:11-jre-slim-buster

# Install GPG for package vefification
RUN apt-get update \
	&& apt-get -y install gnupg wget

# Add the liquibase user and step in the directory
RUN addgroup --gid 1001 liquibase
RUN adduser --disabled-password --uid 1001 --ingroup liquibase liquibase

# Make /liquibase directory and change owner to liquibase
RUN mkdir /liquibase && chown liquibase /liquibase
WORKDIR /liquibase

# Change to the liquibase user
USER liquibase

# Latest Liquibase Release Version
ARG LIQUIBASE_VERSION=4.5.0

# Download, verify, extract
ARG LB_SHA256=4ce45bcbe4660b33eee93e80b9be968631e85566d02d90c0c5306fac8d4dd602
RUN set -x \
  && wget -O liquibase-${LIQUIBASE_VERSION}.tar.gz "https://github.com/liquibase/liquibase/releases/download/v${LIQUIBASE_VERSION}/liquibase-${LIQUIBASE_VERSION}.tar.gz" \
  && sha256sum liquibase-${LIQUIBASE_VERSION}.tar.gz \
  && echo "$LB_SHA256  liquibase-${LIQUIBASE_VERSION}.tar.gz" | sha256sum -c - \
  && tar -xzf liquibase-${LIQUIBASE_VERSION}.tar.gz

# Setup GPG
RUN GNUPGHOME="$(mktemp -d)"

ARG LB_MONGO_VERSION=4.5.0
ARG MDB_JAVA_DRIVER_VERSION=4.3.2
RUN wget -O /liquibase/lib/mongodb.jar https://github.com/liquibase/liquibase-mongodb/releases/download/liquibase-mongodb-${LB_MONGO_VERSION}/liquibase-mongodb-${LB_MONGO_VERSION}.jar
RUN wget -O /liquibase/lib/bson-${MDB_JAVA_DRIVER_VERSION}.jar https://repo1.maven.org/maven2/org/mongodb/bson/${MDB_JAVA_DRIVER_VERSION}/bson-${MDB_JAVA_DRIVER_VERSION}.jar
RUN wget -O /liquibase/lib/mongodb-driver-core-${MDB_JAVA_DRIVER_VERSION}.jar https://repo1.maven.org/maven2/org/mongodb/mongodb-driver-core/${MDB_JAVA_DRIVER_VERSION}/mongodb-driver-core-${MDB_JAVA_DRIVER_VERSION}.jar
RUN wget -O /liquibase/lib/mongodb-driver-sync-${MDB_JAVA_DRIVER_VERSION}.jar https://repo1.maven.org/maven2/org/mongodb/mongodb-driver-sync/${MDB_JAVA_DRIVER_VERSION}/mongodb-driver-sync-${MDB_JAVA_DRIVER_VERSION}.jar

COPY liquibase-docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["liquibase-docker-entrypoint.sh"]
CMD ["--help"]
