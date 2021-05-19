#!/bin/bash



ip netns exec ns1 ip link set veth1 netns 1
ip netns exec ns2 ip link set veth2 netns 1

ip link del veth1 type veth peer name veth2

ip netns exec ns1 ip tuntap del tun0 mode tun
ip netns exec ns2 ip tuntap del tun0 mode tun

ip netns del ns1
ip netns del ns2
