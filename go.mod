module github.com/go-eagle/eagle

go 1.16

require (
	github.com/1024casts/gorm-opentelemetry v1.0.1-0.20210805144709-183269b54068
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/Shopify/sarama v1.19.0
	github.com/Shopify/toxiproxy v2.1.4+incompatible // indirect
	github.com/alicebob/miniredis/v2 v2.15.1
	github.com/dgraph-io/ristretto v0.0.3
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/foolin/gin-template v0.0.0-20190415034731-41efedfb393b
	github.com/frankban/quicktest v1.7.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-contrib/static v0.0.0-20200916080430-d45d9a37d28e
	github.com/gin-gonic/gin v1.7.3
	github.com/go-kratos/aegis v0.1.1
	github.com/go-mail/mail v2.3.1+incompatible
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/extra/redisotel/v8 v8.10.0
	github.com/go-redis/redis/v8 v8.11.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.6
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.3
	github.com/google/uuid v1.2.0
	github.com/gorilla/sessions v1.2.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/consul/api v1.11.0
	github.com/hibiken/asynq v0.23.0
	github.com/jinzhu/gorm v1.9.16
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	github.com/nacos-group/nacos-sdk-go v1.1.1
	github.com/nats-io/nats-server/v2 v2.9.0 // indirect
	github.com/nats-io/nats.go v1.16.1-0.20220906180156-a1017eec10b0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pierrec/lz4 v2.5.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/qiniu/api.v7 v0.0.0-20190520053455-bea02cd22bf4
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.8.0 // indirect
	github.com/spf13/cast v1.4.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.0
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.7.0 // indirect
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/toolkits/net v0.0.0-20160910085801-3f39ab6fe3ce
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	github.com/vearne/gin-timeout v0.1.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/willf/pad v0.0.0-20190207183901-eccfe5d84172
	go.etcd.io/etcd/client/v3 v3.5.1
	go.mongodb.org/mongo-driver v1.5.1
	go.opentelemetry.io/contrib v0.22.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.22.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.22.0
	go.opentelemetry.io/contrib/propagators v0.22.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC3
	go.opentelemetry.io/otel/oteltest v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
	go.opentelemetry.io/otel/trace v1.3.0
	go.uber.org/automaxprocs v1.5.1
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90
	golang.org/x/net v0.0.0-20220107192237-5cfca573fb4d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/tools v0.1.8 // indirect
	google.golang.org/genproto v0.0.0-20220107163113-42d7afdf6368 // indirect
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/mail.v2 v2.3.1 // indirect
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.20.12
)
