kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl add_user k8s k8s
kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl set_user_tags k8s administrator
kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl set_permissions -p / k8s ".*" ".*" ".*"

kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl add_vhost order
kubectl exec --stdin --tty rabbitmq-0 -- rabbitmqctl set_permissions -p order k8s ".*" ".*" ".*"