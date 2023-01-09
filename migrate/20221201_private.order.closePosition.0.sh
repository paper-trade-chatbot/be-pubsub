
kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl eval 'rabbit_exchange:declare({resource, <<"order">>, exchange, <<"order.direct">>}, direct, true, false, false, [],<<"rmqinternal">>).'

kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl eval 'rabbit_amqqueue:declare({resource, <<"order">>, queue, <<"private.order.closePosition.0">>}, true, false, [], none,<<"rmqinternal">>).'

kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl eval 'rabbit_binding:add({binding, {resource, <<"order">>, exchange, <<"order.direct">>}, <<"private.order.closePosition.0">>, {resource, <<"order">>, queue, <<"private.order.closePosition.0">>}, []},<<"rmqinternal">>).'