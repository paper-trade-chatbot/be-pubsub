rabbitmqctl add_user k8s k8s
rabbitmqctl set_user_tags k8s administrator
rabbitmqctl set_permissions -p / k8s ".*" ".*" ".*"

rabbitmqctl add_vhost order
rabbitmqctl set_permissions -p order k8s ".*" ".*" ".*"

rabbitmqctl eval 'rabbit_exchange:declare({resource, <<"order">>, exchange, <<"order.direct">>}, direct, true, false, false, [],<<"rmqinternal">>).'

rabbitmqctl eval 'rabbit_amqqueue:declare({resource, <<"order">>, queue, <<"private.order.openPosition.0">>}, true, false, [], none,<<"rmqinternal">>).'

rabbitmqctl eval 'rabbit_binding:add({binding, {resource, <<"order">>, exchange, <<"order.direct">>}, <<"private.order.openPosition.0">>, {resource, <<"order">>, queue, <<"private.order.openPosition.0">>}, []},<<"rmqinternal">>).'