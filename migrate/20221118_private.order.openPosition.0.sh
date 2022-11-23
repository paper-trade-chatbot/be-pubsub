kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl add_user k8s k8s
kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl set_user_tags k8s administrator
kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl set_permissions -p / k8s ".*" ".*" ".*"

kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl add_vhost order
kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl set_permissions -p order k8s ".*" ".*" ".*"

kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl eval 'rabbit_exchange:declare({resource, <<"order">>, exchange, <<"order.direct">>}, direct, true, false, false, [],<<"rmqinternal">>).'

kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl eval 'rabbit_amqqueue:declare({resource, <<"order">>, queue, <<"private.order.openPosition.0">>}, true, false, [], none,<<"rmqinternal">>).'

kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl eval 'rabbit_binding:add({binding, {resource, <<"order">>, exchange, <<"order.direct">>}, <<"private.order.openPosition.0">>, {resource, <<"order">>, queue, <<"private.order.openPosition.0">>}, []},<<"rmqinternal">>).'