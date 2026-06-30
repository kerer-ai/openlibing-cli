# PR 门禁精选指标报告

> 统计周期: 2026-06-12 ~ 2026-06-30 | 共 87 个仓库 (7 个产品)
>
> 门禁E2E执行列使用 `efficiencyAvgTime*` 字段(不含重试)。

| 产品 | 仓库 | 分支 | 门禁E2E执行P90(min)(不含重试) | 门禁E2E执行平均(min)(不含重试) | P50门禁E2E执行(min)(不含重试) | 构建任务P50(min) | 构建任务P90(min) | 构建任务平均(min) | 构建任务排队P90(min) | 测试任务P90(min) | 测试任务P50(min) | 测试任务平均(min) | 测试任务排队P90(min) | 代码检查任务P90(min) |
|--------------------|------------------|--------------|-----------------|----------------|-----------------|--------------|--------------|-------------|----------------|--------------|--------------|-------------|----------------|----------------|
| Ascend-CANN | AscendNPU-IR-Dev | master | 17.5 | 7.5 | 3.5 | - | - | - | - | - | - | - | - | 2.8 |
| Ascend-CANN | AscendNPU-IR | master | 143.4 | 97.5 | 95.9 | 38.5 | 50.3 | 39.9 | 0.9 | - | - | - | - | 2.8 |
| Ascend-CANN | catlass | master | 71.1 | 39.3 | 43.1 | 2.5 | 3.7 | 2.2 | 0.0 | 67.8 | 41.5 | 35.9 | 0.2 | 2.6 |
| Ascend-CANN | shmem | master | 34.8 | 27.0 | 29.4 | 5.0 | 17.2 | 6.6 | 0.0 | 30.2 | 26.4 | 21.1 | 1.0 | 3.2 |
| Ascend-CANN | ascend-transformer-boost | master | 89.8 | 47.3 | 35.4 | 27.2 | 35.5 | 26.5 | 0.0 | 74.2 | 4.5 | 24.8 | 0.0 | 2.6 |
| Ascend-CANN | catccos | master | 25.3 | 17.4 | 17.1 | - | - | - | - | 15.6 | 15.3 | 14.6 | 0.1 | 2.6 |
| Ascend-CANN | Triton-distributed-ascend | master | 5.5 | 4.9 | 3.4 | - | - | - | - | - | - | - | - | 2.6 |
| Ascend-CANN | sip | master | 14.1 | 8.8 | 7.2 | 2.6 | 2.7 | 2.4 | 0.0 | 9.0 | 6.3 | 6.5 | 0.6 | 2.7 |
| Ascend-CANN | ops-batchinvariant | master | 17.8 | 15.9 | 15.7 | 12.7 | 14.7 | 10.9 | 0.0 | 7.3 | 6.7 | 6.9 | 0.0 | 2.3 |
| Ascend-CANN | ascend-boost-comm | master | 3.4 | 3.3 | 3.3 | 2.4 | 2.4 | 2.4 | 0.0 | 1.9 | 1.8 | 1.8 | 0.0 | 2.5 |
| Ascend-CANN | ops-collections | master | 4.6 | 4.6 | 4.6 | - | - | - | - | 2.4 | 2.4 | 2.4 | 0.0 | 3.6 |
| FrameworkPTAdapter | pytorch | master | 76.2 | 40.5 | 35.6 | 9.3 | 16.0 | 10.4 | 0.1 | 73.0 | 30.4 | 38.5 | 0.1 | 3.1 |
| FrameworkPTAdapter | op-plugin | 7.3.0 | 87.7 | 32.7 | 22.1 | 12.4 | 21.8 | 16.2 | 0.1 | 75.5 | 37.6 | 43.9 | 0.1 | 2.9 |
| FrameworkPTAdapter | torchair | master | 14.6 | 9.8 | 12.9 | 3.8 | 7.6 | 4.5 | 0.1 | 14.5 | 12.6 | 12.8 | 0.1 | 3.0 |
| FrameworkPTAdapter | apex | master | 4.4 | 4.4 | 4.4 | 1.8 | 1.8 | 1.8 | 0.0 | - | - | - | - | 2.5 |
| MindCluster | mind-cluster | master | 24.2 | 12.5 | 6.3 | 5.0 | 24.4 | 11.6 | 0.2 | 13.4 | 4.8 | 6.8 | 0.3 | 2.3 |
| MindCluster | RecSDK | develop | 25.1 | 14.9 | 14.1 | 4.3 | 6.7 | 4.5 | 0.2 | 14.3 | 8.2 | 8.3 | 0.3 | 2.3 |
| MindCluster | ascend-deployer | dev | 4.9 | 4.3 | 4.2 | 0.9 | 1.1 | 0.9 | 0.2 | 1.2 | 1.0 | 1.0 | 0.2 | 2.4 |
| MindCluster | fbgemm-ascend | main | 21.5 | 10.2 | 8.0 | 6.6 | 20.1 | 9.7 | 0.1 | - | - | - | - | 2.2 |
| MindCluster | AgentSDK | master | 28.0 | 14.8 | 13.1 | 2.6 | 8.5 | 3.8 | 0.2 | 18.5 | 11.2 | 12.1 | 0.2 | 2.2 |
| MindCluster | MindCluster-AscendNPUBurn | master | 14.2 | 10.4 | 8.3 | 7.3 | 13.1 | 9.4 | 0.2 | 1.1 | 1.0 | 1.0 | 0.1 | 2.3 |
| MindCluster | IndexSDK | master | 28.9 | 16.6 | 13.2 | 6.4 | 26.9 | 9.7 | 0.1 | 23.9 | 11.5 | 13.7 | 0.2 | 2.3 |
| MindCluster | perf-reference-ascend | master | 3.7 | 3.0 | 3.1 | - | - | - | - | - | - | - | - | 2.5 |
| MindCluster | RAGSDK | master | 25.9 | 15.0 | 16.0 | 4.7 | 9.0 | 4.9 | 0.1 | 33.2 | 16.4 | 19.7 | 0.2 | 2.3 |
| MindCluster | MultimodalSDK | master | 30.5 | 18.8 | 17.6 | 6.1 | 18.8 | 8.8 | 0.1 | 31.2 | 16.3 | 18.8 | 0.2 | 2.2 |
| MindCluster | faiss | main | 6.8 | 4.9 | 5.2 | 5.6 | 7.5 | 6.0 | 0.1 | - | - | - | - | 2.2 |
| MindCluster | mind-cluster_for_lingqu | master | 17.9 | 8.7 | 6.4 | 4.3 | 17.1 | 7.1 | 0.1 | 8.1 | 3.2 | 4.4 | 0.1 | 2.6 |
| MindCluster | VisionSDK | master | 47.3 | 34.9 | 34.2 | 18.3 | 39.9 | 22.9 | 1.0 | 41.1 | 33.5 | 30.7 | 0.8 | 2.4 |
| MindCluster | HierarchicalKV-ascend | develop | 4.1 | 3.4 | 3.3 | 1.1 | 1.9 | 1.3 | 0.0 | - | - | - | - | 2.2 |
| MindCluster | mindcluster-deploy | master | 3.1 | 2.4 | 3.0 | - | - | - | - | - | - | - | - | 2.1 |
| MindCluster | mindsdk-referenceapps | master | 2.1 | 1.9 | 1.8 | 0.9 | 0.9 | 0.9 | 0.1 | - | - | - | - | 2.5 |
| MindCluster | OMSDK | master | 11.8 | 10.7 | 10.7 | 9.6 | 12.8 | 9.6 | 0.1 | 16.5 | 12.6 | 12.6 | 0.1 | 2.2 |
| MindCluster | text-embeddings-inference | main | 17.6 | 16.4 | 16.4 | - | - | - | - | 16.6 | 15.3 | 15.3 | 0.1 | 2.1 |
| MindCluster | MEF | master | 12.1 | 12.1 | 12.1 | 11.2 | 11.2 | 11.2 | 0.0 | 7.4 | 7.4 | 7.4 | 0.1 | 2.2 |
| MindIE | vllm-ascend | main | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |
| MindIE | MindIE-PyMotor | master | 10.9 | 7.6 | 7.9 | 2.8 | 4.7 | 3.1 | 0.1 | 5.8 | 3.0 | 3.4 | 0.1 | 2.3 |
| MindIE | MindIE-SD | master | 10.5 | 8.2 | 7.7 | 5.5 | 8.3 | 6.4 | 0.5 | 3.6 | 2.5 | 3.5 | 0.1 | 2.3 |
| MindIE | MindIE-LLM | master | 59.1 | 27.9 | 23.8 | 17.9 | 40.4 | 19.4 | 0.1 | 26.1 | 21.2 | 16.3 | 0.1 | 2.4 |
| MindIE | MindIE-Motor | master | 123.2 | 50.2 | 20.5 | 5.5 | 74.8 | 27.8 | 0.1 | 132.6 | 23.8 | 54.5 | 0.1 | 2.3 |
| MindSpeed | MindSpeed-LLM | master | 62.9 | 16.6 | 3.8 | - | - | - | - | 60.5 | 29.0 | 28.5 | 0.0 | 2.4 |
| MindSpeed | MindSpeed-MM | master | 42.1 | 20.0 | 15.6 | 1.3 | 2.0 | 1.4 | 0.1 | 28.3 | 22.5 | 18.1 | 0.0 | 2.3 |
| MindSpeed | MindSpeed | master | 43.2 | 22.3 | 23.0 | 1.0 | 1.9 | 1.2 | 0.0 | 63.2 | 35.1 | 37.3 | 21.0 | 3.0 |
| MindSpeed | DrivingSDK | master | 12.1 | 6.4 | 4.7 | 6.5 | 9.0 | 6.4 | 0.1 | 2.1 | 0.3 | 0.9 | 0.0 | 2.9 |
| MindSpeed | ModelZoo-PyTorch | master | 3.5 | 3.7 | 2.7 | - | - | - | - | - | - | - | - | 3.0 |
| MindSpeed | MindSpeed-Ops | master | 38.6 | 27.7 | 27.7 | - | - | - | - | 61.5 | 31.7 | 36.4 | 2.2 | 2.3 |
| MindSpeed | TransformerEngineNPU | main | 13.3 | 8.1 | 5.5 | - | - | - | - | 17.7 | 8.2 | 9.5 | 2.5 | 2.4 |
| MindSpeed | MegatronAdaptor | main | 30.0 | 11.6 | 5.1 | - | - | - | - | 35.9 | 2.2 | 14.9 | 0.8 | 7.0 |
| MindSpeed | MindSpeed-Bridge | master | 15.0 | 9.6 | 10.4 | - | - | - | - | 14.2 | 9.3 | 8.9 | 0.0 | 2.4 |
| MindSpeed | slime-ascend | main | 24.1 | 13.5 | 18.2 | - | - | - | - | 29.6 | 25.2 | 24.5 | 0.0 | 2.1 |
| MindSpeed | modelzoo-GPL | master | 2.3 | 2.1 | 2.1 | - | - | - | - | - | - | - | - | 2.3 |
| MindSpeed | MindSpeed-RL | master | 0.4 | 0.4 | 0.4 | - | - | - | - | - | - | - | - | - |
| MindSpeed | apex | master | 4.4 | 4.4 | 4.4 | 1.8 | 1.8 | 1.8 | 0.0 | - | - | - | - | 2.5 |
| MindSpore | mindformers | master | 31.9 | 18.6 | 17.3 | - | - | - | - | 56.5 | 32.7 | 34.7 | 8.4 | 2.6 |
| MindStudio | msmodeling | master | 24.1 | 10.2 | 6.5 | 1.6 | 1.9 | 1.6 | 0.1 | 26.1 | 3.9 | 10.6 | 2.6 | 2.5 |
| MindStudio | msmodelslim | master | 29.5 | 14.3 | 14.8 | 1.5 | 2.8 | 1.8 | 0.0 | 22.5 | 18.9 | 19.3 | 5.3 | 2.5 |
| MindStudio | msinsight | master | 35.0 | 27.6 | 33.3 | 12.2 | 13.6 | 12.7 | 1.5 | 18.0 | 16.4 | 16.5 | 0.0 | 2.6 |
| MindStudio | msprobe | master | 13.1 | 8.5 | 8.2 | 0.9 | 4.0 | 2.3 | 0.1 | 7.0 | 5.8 | 6.0 | 0.5 | 2.4 |
| MindStudio | msprof | master | 17.4 | 7.6 | 1.9 | 3.5 | 5.6 | 4.3 | 0.0 | 10.7 | 9.1 | 9.3 | 0.8 | 2.5 |
| MindStudio | mssanitizer | master | 41.8 | 23.6 | 19.8 | 14.2 | 17.9 | 14.1 | 0.0 | 17.9 | 4.1 | 9.1 | 15.5 | 2.5 |
| MindStudio | msserviceprofiler | master | 11.3 | 6.8 | 6.8 | 2.9 | 3.6 | 3.9 | 0.0 | 8.2 | 2.8 | 3.9 | 0.0 | 2.5 |
| MindStudio | msagent | master | 4.7 | 3.5 | 4.3 | 0.8 | 1.0 | 0.8 | 0.0 | - | - | - | - | 2.3 |
| MindStudio | msopcom | master | 22.6 | 19.0 | 12.1 | 2.2 | 2.4 | 2.1 | 4.5 | 5.8 | 5.5 | 5.5 | 2.0 | 2.5 |
| MindStudio | mspti | master | 12.3 | 6.0 | 4.5 | 2.1 | 2.9 | 2.2 | 0.1 | 7.0 | 6.2 | 6.2 | 1.0 | 2.4 |
| MindStudio | msdebug | master | 40.2 | 26.1 | 25.4 | 8.1 | 12.0 | 10.1 | 0.0 | 12.5 | 9.9 | 10.8 | 0.9 | 9.8 |
| MindStudio | msmemscope | master | 30.9 | 17.4 | 17.5 | 6.8 | 15.1 | 9.8 | 0.0 | 16.4 | 14.3 | 14.8 | 9.2 | 2.6 |
| MindStudio | msopprof | master | 24.1 | 10.9 | 1.6 | 9.4 | 10.2 | 9.0 | 0.1 | 10.7 | 9.5 | 9.7 | 4.9 | 2.5 |
| MindStudio | msmonitor | master | 9.4 | 5.0 | 2.4 | 3.2 | 7.0 | 4.1 | 0.1 | 9.8 | 9.1 | 8.7 | 0.9 | 2.2 |
| MindStudio | msit | master | 4.0 | 2.8 | 1.7 | - | - | - | - | 2.3 | 1.3 | 1.3 | 0.1 | 2.3 |
| MindStudio | msprof-analyze | master | 12.3 | 5.6 | 2.4 | 1.6 | 1.9 | 1.7 | 0.0 | 10.3 | 8.9 | 9.0 | 5.6 | 2.5 |
| MindStudio | mskpp | master | 15.6 | 7.6 | 1.6 | 2.5 | 3.0 | 2.6 | 0.0 | 30.8 | 3.7 | 13.2 | 0.0 | 2.5 |
| MindStudio | mstt | master | 2.9 | 2.0 | 2.1 | - | - | - | - | 2.8 | 2.2 | 2.0 | 0.1 | 2.4 |
| MindStudio | msboost | master | 3.3 | 2.6 | 3.0 | - | - | - | - | - | - | - | - | 2.6 |
| MindStudio | msot | master | 12.4 | 5.0 | 1.4 | 19.4 | 19.4 | 19.4 | 0.0 | - | - | - | - | 2.4 |
| MindStudio | mstx | master | 9.1 | 5.2 | 5.0 | 1.9 | 2.2 | 1.9 | 0.1 | 1.0 | 0.9 | 0.9 | 3.2 | 2.3 |
| MindStudio | mskl | master | 6.9 | 5.5 | 6.5 | 1.8 | 1.9 | 1.8 | 0.1 | 2.0 | 0.8 | 1.2 | 0.0 | 2.2 |
| MindStudio | msopgen | master | 7.1 | 6.2 | 6.1 | 1.9 | 2.1 | 1.9 | 0.0 | 2.3 | 1.6 | 1.7 | 0.0 | 2.3 |
| MindStudio | msmodelslim_for_lingqu | master | 31.1 | 23.5 | 27.2 | 0.9 | 1.1 | 1.0 | 0.0 | 26.2 | 20.2 | 21.4 | 2.2 | 2.3 |
| MindStudio | msoptuner | master | 10.7 | 9.9 | 9.9 | 4.5 | 4.5 | 4.5 | 0.1 | 2.5 | 2.3 | 2.3 | 0.0 | 2.3 |
| MindStudio | msprof_for_lingqu | master | 15.6 | 10.9 | 10.9 | 2.7 | 3.1 | 2.7 | 0.1 | 5.5 | 3.7 | 3.7 | 7.9 | 2.4 |
| MindStudio | mscommreport | master | 0.7 | 0.7 | 0.7 | - | - | - | - | - | - | - | - | - |
| MindStudio | msinsight_for_lingqu | master | 25.7 | 25.7 | 25.7 | 9.8 | 9.8 | 9.8 | 6.9 | 10.7 | 10.7 | 10.7 | 0.0 | 3.8 |
| MindStudio | msmonitor_for_lingqu | master | 5.0 | 5.0 | 5.0 | 2.1 | 2.1 | 2.1 | 0.0 | - | - | - | - | 2.2 |
| MindStudio | msopcom_for_lingqu | master | 3.7 | 3.7 | 3.7 | 0.7 | 0.7 | 0.7 | 0.0 | - | - | - | - | 2.3 |
| MindStudio | msopmodeling | master | 8.3 | 8.3 | 8.3 | 0.7 | 0.7 | 0.7 | 0.0 | 6.8 | 6.8 | 6.8 | 0.0 | 2.3 |
| MindStudio | msopprof_for_lingqu | master | 4.5 | 4.5 | 4.5 | 0.8 | 0.8 | 0.8 | 0.0 | - | - | - | - | 2.3 |
| MindStudio | msprof-analyze_for_lingqu | master | 6.4 | 6.4 | 6.4 | 0.5 | 0.5 | 0.5 | 0.0 | 2.7 | 2.7 | 2.7 | 0.0 | 2.5 |
| MindStudio | mspti_for_lingqu | master | 3.0 | 3.0 | 3.0 | - | - | - | - | - | - | - | - | 2.2 |

## 概览统计

| 产品 | 仓库数 |
|------|--------|
| Ascend-CANN | 11 |
| FrameworkPTAdapter | 4 |
| MindCluster | 19 |
| MindIE | 5 |
| MindSpeed | 13 |
| MindSpore | 1 |
| MindStudio | 34 |

---
*报告生成时间: 2026-06-30 15:39:29 | 数据来源: openlibing-ops repo-pr-pipeline (efficiencyAvgTime* 字段)*
