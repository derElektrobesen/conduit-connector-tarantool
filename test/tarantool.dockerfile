FROM docker.io/tarantool/tarantool:2.11-centos7

RUN yum-config-manager --disable tarantool_2,tarantool_2-source,tarantool-modules,tarantool_modules-source
RUN yum makecache

RUN yum -y groupinstall "Development Tools"
RUN yum -y install cmake

COPY tarantool/customer-scm-1.rockspec /tmp/
RUN tarantoolctl rocks --tree .rocks install --only-deps /tmp/customer-scm-1.rockspec
