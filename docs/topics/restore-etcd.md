# Scenario:
The Kubernetes cluster consists of 3 control plane nodes. One of these nodes, specifically named (for example, *-master-*-0), has encountered a corruption issue with its etcd data disk. 

The steps outlined below detail the recovery procedure for this situation. 

<span style="font-size:30px;color:red;">**It is important to note that during the recovery process, the cluster will be temporarily unavailable.**</span>

# 1. Take a snapshot of the etcd cluster on a **healthy** node.
 The script assume the **healthy** node is the second master node. (\*-master-\*-1)

ssh into the second control plane node.

<span style="font-size:20px;color:yellow;">**Please replace the FQDN (akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com) of the public ip to your own environment.**</span>

```
ssh -p 2201 azureuser@akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com
```

```
sudo su
export ETCDCTL_API=3
export ETCDCTL_CERT=/etc/kubernetes/certs/etcdclient.crt
export ETCDCTL_CACERT=/etc/kubernetes/certs/ca.crt
export ETCDCTL_KEY=/etc/kubernetes/certs/etcdclient.key
etcdctl snapshot save /home/azureuser/snapshot.db
chmod 644 /home/azureuser/snapshot.db
```
# 2. Detach the etcd data disk from portal.
## 2.1 Find the control plane node with corrupted etcd disk. (for example *-master-*-0) on portal.
## 2.2 Stop the VM if it is still running
### 2.2.1 Click the Refresh and Wait for the VM status from Deallocating to Stopped (deallocated)
## 2.3 Go to Disks tab
## 2.4 Noted down the disk information attached on LUN 0
    - Name
    - disk size
    - storage account type
    - tags
## 2.5 Detch the etcd disk and click "Save"
## 2.6 Wait for the operation completion.
# 3. Find and delete the the etcd data disk from portal.
# 4. Create and Attach the new created data disk to control plane VM 
## 4.1 Find the control plane node with corrupted etcd disk. (for example *-master-*-0) on portal.
## 4.2 Go to Disks tab
## 4.3 Click "Add data disk" and esure that the **LUN number is 0**
## 4.4 Choose "Create disk" in the dropdown list
## 4.5 Create new empty data disk with same name, disk size, resource group, storage account type and tags.
## 4.6 Wait for the disk creation to complete
## 4.7 Click Save to update the control plane VM.
## 4.8 Start the control plane VM

# 5. Copy the snapshot to all other nodes to location /home/azureuser/snapshot.db

## 5.1 copy the snapshot from the second master node to the jumpbox
<span style="font-size:20px;color:yellow;">**Please replace the FQDN (akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com) of the public ip to your own environment.**</span>
```
scp -P 2201 azureuser@akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com:/home/azureuser/snapshot.db .
```

## 5.2 copy the snapshot into the first and third master
<span style="font-size:20px;color:yellow;">**Please replace the FQDN (akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com) of the public ip to your own environment.**</span>
```
scp -P 22 ./snapshot.db azureuser@akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com:/home/azureuser
scp -P 2202 ./snapshot.db azureuser@akse0007.redmond.cloudapp.ext-n25r1304.masd.stbtest.microsoft.com:/home/azureuser
```
# 6. Disable and Stop the etcd service on all threes masters node
Open ssh session to all three master node 

```
sudo su
```
```
systemctl disable etcd
```
```
systemctl stop etcd
```

# 7. Set the environment variables
Reuse the SSH session in previous step


## 7.1 Set environment variables for master

<span style="font-size:20px;color:yellow;">**Please replace the FIRST_MASTER_IP,SECOND_MASTER_IP,THIRD_MASTER_IP with the value in your environment.**</span>

```
export FIRST_MASTER_IP=10.240.255.5
export SECOND_MASTER_IP=10.240.255.6
export THIRD_MASTER_IP=10.240.255.7
```
## 7.2 set environment variables etcdctl command
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
# 8 Restore the etcd snap shot on all the master nodes
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
# 9. Start the etcd service all all master node and verify the service status
```
systemctl start etcd
sleep 10
```
```
systemctl status etcd
```
# 10. Verify the the cluster status 

```
etcdctl endpoint status --write-out=table
```

```
etcdctl endpoint health --write-out=table
```

# 11 Enable the etcd service
```
systemctl enable etcd
```