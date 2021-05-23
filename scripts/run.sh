ip netns exec ns1 ./svpn config/server.yaml
ip netns exec ns2 ./svpn config/client.yml

ip netns exec ns2 ping 6.6.6.100
