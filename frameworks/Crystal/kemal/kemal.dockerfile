FROM crystallang/crystal:0.26.1

WORKDIR /kemal
COPY views views
COPY run.sh run.sh
COPY server-postgres.cr server-postgres.cr
COPY shard.lock shard.lock
COPY shard.yml shard.yml

ENV GC_MARKERS 1
ENV KEMAL_ENV production
ENV DATABASE_URL postgres://benchmarkdbuser:benchmarkdbpass@tfb-database:5432/hello_world?initial_pool_size=8&max_pool_size=8&max_idle_pool_size=8

RUN shards install
RUN crystal build --release --no-debug server-postgres.cr

CMD bash run.sh