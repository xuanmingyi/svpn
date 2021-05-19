#!/bin/bash

ip netns add ns1
ip netns add ns2

ip link add veth1 type veth peer name veth2

ip link set veth1 netns ns1
ip link set veth2 netns ns2

ip netns exec ns1 ip addr add 18.1.1.101/24 dev veth1
ip netns exec ns1 ip link set veth1 up
ip netns exec ns2 ip addr add 18.1.1.102/24 dev veth2
ip netns exec ns2 ip link set veth2 up

ip netns exec ns1 ip tuntap add tun0 mode tun
ip netns exec ns1 ip addr add 6.6.6.101/24 dev tun0
ip netns exec ns1 ip link set tun0 up
ip netns exec ns1 ip link set lo up

ip netns exec ns2 ip tuntap add tun0 mode tun
ip netns exec ns2 ip addr add 6.6.6.102/24 dev tun0
ip netns exec ns2 ip link set tun0 up
ip netns exec ns2 ip link set lo up
