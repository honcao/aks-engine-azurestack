

# 1. Take a snapshot of the etcd cluster on a healthy node.
 The script assume the healthy node is the second master node. (*-master-*1)

```
sudo su
export ETCDCTL_CERT=/etc/kubernetes/certs/etcdclient.crt
export ETCDCTL_CACERT=/etc/kubernetes/certs/ca.crt
export ETCDCTL_KEY=/etc/kubernetes/certs/etcdclient.key
etcdctl snapshot save /home/azureuser/snapshot.db
chmod 644 /home/azureuser/snapshot.db
```
# 2. Copy the snapshot to all other nodes to location /home/azureuser/snapshot.db

## 2.1 copy the snapshot from the second master node to the jumpbox
**Please replace the FQDN (akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com) of the public ip to your own environment.**
```
scp -P 2201 azureuser@akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com:/home/azureuser/snapshot.db .
```

## 2.2 copy the snapshot into the first and third master
**Please replace the FQDN (akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com) of the public ip to your own environment.**
```
scp -P 22 ./snapshot.db azureuser@akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com:/home/azureuser
scp -P 2202 ./snapshot.db azureuser@akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com:/home/azureuser
```
# 3. Disable and Stop the etcd service on all threes masters node
Open ssh session to all three master node 
```
sudo su
systemctl disable etcd
systemctl stop etcd
```

# 4. Restore the etcd snap shot on all the master nodes
Reuse the SSH session in previous step


## 4.1 Set environment variables for master
**Please replace the FIRST_MASTER_IP,SECOND_MASTER_IP,THIRD_MASTER_IP with the value in your environment.**
```
export FIRST_MASTER_IP=10.240.255.5
export SECOND_MASTER_IP=10.240.255.6
export THIRD_MASTER_IP=10.240.255.7
```
## 4.2 set environment variables etcdctl command
```
export ETCDCTL_CERT=/etc/kubernetes/certs/etcdclient.crt
export ETCDCTL_CACERT=/etc/kubernetes/certs/ca.crt
export ETCDCTL_KEY=/etc/kubernetes/certs/etcdclient.key
export ETCDCTL_API=3
export ETCD_ARGUMENTS=$(cat /etc/default/etcd)
export INITIAL_CLUSTER=$(echo $ETCD_ARGUMENTS | cut -d' ' -f22)
echo "INITIAL_CLUSTER $INITIAL_CLUSTER"
export PEER_URL=$(echo $ETCD_ARGUMENTS | cut -d' ' -f16)
echo "PEER_URL $PEER_URL"
export ETCDCTL_ENDPOINTS=https://${FIRST_MASTER_IP}:2379,https://${SECOND_MASTER_IP}:2379,https://${THIRD_MASTER_IP}:2379
echo "ETCDCTL_ENDPOINTS $ETCDCTL_ENDPOINTS"
export ETCD_NAME=$(hostname)
echo "ETCD_NAME $ETCD_NAME"
export HOST_IP=$(hostname -i)
echo "HOST_IP $HOST_IP"
export ETCD_ADVERTISE_PEER_URLS="https://${HOST_IP}:2380"
echo "ETCD_ADVERTISE_PEER_URLS $ETCD_ADVERTISE_PEER_URLS"

```
## 4.2 restore the etcd database
```
# Remove existing memeber data folder
rm -r /var/lib/etcddisk/member
# Restore the database to a temp folder
etcdctl snapshot restore /home/azureuser/snapshot.db  --name $ETCD_NAME --initial-cluster $INITIAL_CLUSTER   --initial-cluster-token k8s-etcd-cluster   --initial-advertise-peer-urls $ETCD_ADVERTISE_PEER_URLS  --data-dir /var/lib/etcddisk/temp

# Copy the database from the temp folder
cp -r /var/lib/etcddisk/temp/member /var/lib/etcddisk/
# Remove the temp folder
rm -r /var/lib/etcddisk/temp

# Change the owner to etcd
chown -R etcd:etcd /var/lib/etcddisk
```
# 5. Start the etcd service all all master node and verify the service status
```
systemctl start etcd
sleep 10
systemctl status etcd
```
# 6. Verify the the cluster status 

```
etcdctl endpoint status --write-out=table
```

```
etcdctl endpoint health --write-out=table
```