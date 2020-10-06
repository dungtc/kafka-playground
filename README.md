# 1. Messaging with Simple Publish Subscribe (pub/sub)
A publisher or multiple publisher write messages to many partitions and single or multiple consumer groups consume those messages.

**Consumer group** is a group of consumers, it can parallelise the data consumption.

When a new consumer joins to consumer group or remove a consumer, **rebalacing** happen

**Consumer group coordinator** implement rebalacing strategy and manage the state of the group. It auto rebalace message to appropriate consumer in each consumer group

![alt text](https://images.squarespace-cdn.com/content/v1/56894e581c1210fead06f878/1512813188413-EEI12VI1FMQLJ4XMRQTJ/ke17ZwdGBToddI8pDm48kJMyHperWZxre3bsQoFNoPhZw-zPPgdn4jUwVcJE1ZvWEtT5uBSRWt4vQZAgTJucoTqqXjS3CfNDSuuf31e0tVGDclntk9GVn4cF1XFdv7wlNvED_LyEM5kIdmOo2jMRZpu3E9Ef3XsXP1C_826c-iU/KafkaPubSub.png?format=500w)
