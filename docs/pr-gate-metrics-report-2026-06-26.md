# PR 门禁精选指标报告

> 统计周期: 2026-06-12 ~ 2026-06-26 | 共 86 个仓库 (7 个产品)
>
> **✅ 已修正:** 门禁E2E执行列使用 `efficiencyAvgTime*` 字段(不含重试)，已于 2026-06-30 修正。

| 产品 | 仓库 | 分支 | 门禁E2E执行P90(min)(不含重试) | 门禁E2E执行平均(min)(不含重试) | P50门禁E2E执行(min)(不含重试) | 构建任务P50(min) | 构建任务P90(min) | 构建任务平均(min) | 构建任务排队P90(min) | 测试任务P90(min) | 测试任务P50(min) | 测试任务平均(min) | 测试任务排队P90(min) | 代码检查任务P90(min) |
|--------------------|------------------|--------------|-----------------|----------------|-----------------|--------------|--------------|-------------|----------------|--------------|--------------|-------------|----------------|----------------|
| Ascend-CANN | AscendNPU-IR-Dev | master | 18.3 | 7.9 | 3.6 | - | - | - | - | - | - | - | - | 2.8 |
| Ascend-CANN | catlass | master | 71.2 | 36.9 | 35.2 | 2.4 | 3.7 | 2.1 | 0.0 | 67.5 | 31.8 | 32.9 | 0.3 | 2.6 |
| Ascend-CANN | AscendNPU-IR | master | 154.0 | 104.2 | 104.7 | 37.3 | 49.8 | 38.1 | 1.6 | - | - | - | - | 3.0 |
| Ascend-CANN | shmem | master | 36.1 | 28.2 | 29.3 | 1.4 | 15.2 | 4.8 | 0.0 | 30.9 | 26.5 | 23.4 | 1.6 | 3.1 |
| Ascend-CANN | ascend-transformer-boost | master | 105.8 | 54.5 | 37.2 | 25.7 | 27.6 | 24.1 | 0.0 | 78.0 | 6.2 | 32.0 | 0.0 | 2.6 |
| Ascend-CANN | Triton-distributed-ascend | master | 5.5 | 4.9 | 3.4 | - | - | - | - | - | - | - | - | 2.6 |
| Ascend-CANN | catccos | master | 27.4 | 19.0 | 17.5 | - | - | - | - | 15.3 | 14.7 | 14.0 | 0.1 | 3.0 |
| Ascend-CANN | sip | master | 15.0 | 10.6 | 10.0 | 2.4 | 2.6 | 2.3 | 0.0 | 9.2 | 7.8 | 7.6 | 0.7 | 2.5 |
| Ascend-CANN | ops-batchinvariant | master | 17.8 | 15.9 | 15.7 | 12.7 | 14.7 | 10.9 | 0.0 | 7.3 | 6.7 | 6.9 | 0.0 | 2.3 |
| Ascend-CANN | ascend-boost-comm | master | 3.4 | 3.4 | 3.4 | 2.3 | 2.3 | 2.3 | 0.0 | 1.7 | 1.7 | 1.7 | 0.0 | 2.5 |
| Ascend-CANN | ops-collections | master | 4.6 | 4.6 | 4.6 | - | - | - | - | 2.4 | 2.4 | 2.4 | 0.0 | 3.6 |
| FrameworkPTAdapter | pytorch | master | 80.8 | 44.1 | 40.7 | 9.4 | 16.5 | 10.8 | 0.1 | 74.4 | 32.6 | 40.2 | 0.1 | 3.2 |
| FrameworkPTAdapter | op-plugin | 7.3.0 | 82.3 | 30.4 | 21.9 | 12.6 | 23.2 | 16.9 | 0.1 | 74.0 | 35.9 | 39.4 | 0.1 | 2.9 |
| FrameworkPTAdapter | torchair | master | 15.2 | 10.6 | 13.1 | 3.9 | 8.3 | 4.7 | 0.1 | 14.7 | 12.5 | 12.8 | 0.1 | 3.5 |
| FrameworkPTAdapter | apex | master | 4.4 | 4.4 | 4.4 | 1.8 | 1.8 | 1.8 | 0.0 | - | - | - | - | 2.5 |
| MindCluster | mind-cluster | master | 33.5 | 15.0 | 7.3 | 5.5 | 28.8 | 13.9 | 0.2 | 16.7 | 5.3 | 7.7 | 0.3 | 2.3 |
| MindCluster | RecSDK | develop | 25.5 | 14.6 | 13.1 | 4.2 | 6.4 | 4.4 | 0.2 | 13.0 | 8.2 | 8.2 | 0.8 | 2.3 |
| MindCluster | ascend-deployer | dev | 4.5 | 4.4 | 4.2 | 0.9 | 1.1 | 0.9 | 0.1 | 1.3 | 1.0 | 1.1 | 0.2 | 2.3 |
| MindCluster | fbgemm-ascend | main | 22.1 | 10.6 | 8.1 | 6.8 | 21.1 | 9.9 | 0.1 | - | - | - | - | 2.2 |
| MindCluster | AgentSDK | master | 30.6 | 15.4 | 13.5 | 2.7 | 8.6 | 3.9 | 0.2 | 21.6 | 11.3 | 12.2 | 0.2 | 2.2 |
| MindCluster | IndexSDK | master | 28.9 | 17.5 | 15.0 | 6.6 | 27.0 | 9.9 | 0.1 | 23.9 | 11.3 | 13.8 | 0.2 | 2.3 |
| MindCluster | MindCluster-AscendNPUBurn | master | 16.1 | 10.5 | 8.4 | 7.3 | 15.2 | 9.5 | 0.2 | 1.2 | 1.0 | 1.1 | 0.1 | 2.3 |
| MindCluster | perf-reference-ascend | master | 3.7 | 3.2 | 3.3 | - | - | - | - | - | - | - | - | 2.6 |
| MindCluster | MultimodalSDK | None | 32.1 | 20.8 | 17.7 | 6.1 | 18.8 | 8.8 | 0.1 | 31.2 | 16.3 | 18.8 | 0.2 | 2.2 |
| MindCluster | RAGSDK | master | 27.2 | 16.6 | 16.9 | 4.8 | 8.3 | 5.2 | 0.1 | 36.1 | 18.0 | 21.9 | 0.2 | 2.2 |
| MindCluster | faiss | main | 6.8 | 5.1 | 5.3 | 5.8 | 7.7 | 6.1 | 0.1 | - | - | - | - | 2.2 |
| MindCluster | mind-cluster_for_lingqu | master | 18.9 | 9.2 | 6.4 | 4.3 | 17.9 | 7.5 | 0.1 | 8.7 | 3.2 | 4.6 | 0.1 | 2.6 |
| MindCluster | HierarchicalKV-ascend | develop | 4.1 | 3.4 | 3.3 | 1.1 | 1.9 | 1.3 | 0.0 | - | - | - | - | 2.2 |
| MindCluster | VisionSDK | master | 50.2 | 32.8 | 27.2 | 18.3 | 39.9 | 22.9 | 1.0 | 35.1 | 30.7 | 25.7 | 1.4 | 2.4 |
| MindCluster | mindsdk-referenceapps | master | 2.1 | 1.9 | 1.8 | 0.9 | 0.9 | 0.9 | 0.1 | - | - | - | - | 2.5 |
| MindCluster | OMSDK | master | 11.8 | 10.7 | 10.7 | 9.6 | 12.8 | 9.6 | 0.1 | 16.5 | 12.6 | 12.6 | 0.1 | 2.2 |
| MindCluster | text-embeddings-inference | main | 17.6 | 16.4 | 16.4 | - | - | - | - | 16.6 | 15.3 | 15.3 | 0.1 | 2.1 |
| MindCluster | MEF | master | 12.1 | 12.1 | 12.1 | 11.2 | 11.2 | 11.2 | 0.0 | 7.4 | 7.4 | 7.4 | 0.1 | 2.2 |
| MindCluster | mindcluster-deploy | master | 3.0 | 3.0 | 3.0 | - | - | - | - | - | - | - | - | 1.9 |
| MindIE | vllm-ascend | main | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |
| MindIE | MindIE-PyMotor | master | 11.0 | 7.7 | 8.0 | 2.8 | 4.7 | 3.1 | 0.1 | 6.0 | 3.0 | 3.5 | 0.1 | 2.3 |
| MindIE | MindIE-SD | master | 11.4 | 8.8 | 7.9 | 5.5 | 12.1 | 7.0 | 1.4 | 3.0 | 2.4 | 3.9 | 0.1 | 2.3 |
| MindIE | MindIE-LLM | master | 55.2 | 25.4 | 18.4 | 17.5 | 41.0 | 18.0 | 0.1 | 24.9 | 12.4 | 14.6 | 0.1 | 2.4 |
| MindIE | MindIE-Motor | master | 123.2 | 50.2 | 20.5 | 5.5 | 74.8 | 27.8 | 0.1 | 132.6 | 23.8 | 54.5 | 0.1 | 2.3 |
| MindSpeed | MindSpeed-LLM | master | 63.0 | 16.8 | 3.9 | - | - | - | - | 60.5 | 28.0 | 27.3 | 0.0 | 2.4 |
| MindSpeed | MindSpeed-MM | master | 41.9 | 19.8 | 15.6 | 1.3 | 2.0 | 1.4 | 0.1 | 28.3 | 19.7 | 17.5 | 0.0 | 2.4 |
| MindSpeed | MindSpeed | master | 44.3 | 24.0 | 24.8 | 1.0 | 1.8 | 1.2 | 0.0 | 63.9 | 39.0 | 37.9 | 21.3 | 3.1 |
| MindSpeed | DrivingSDK | master | 10.8 | 5.8 | 3.9 | 4.1 | 10.7 | 6.2 | 0.1 | 1.8 | 0.3 | 0.9 | 0.0 | 3.3 |
| MindSpeed | MindSpeed-Ops | master | 38.7 | 30.0 | 33.4 | - | - | - | - | 61.5 | 31.7 | 36.4 | 2.2 | 2.3 |
| MindSpeed | ModelZoo-PyTorch | master | 3.6 | 3.9 | 2.7 | - | - | - | - | - | - | - | - | 3.0 |
| MindSpeed | TransformerEngineNPU | main | 15.2 | 8.1 | 5.5 | - | - | - | - | 17.7 | 8.2 | 9.5 | 2.7 | 2.6 |
| MindSpeed | MegatronAdaptor | main | 30.7 | 13.0 | 5.4 | - | - | - | - | 35.9 | 2.2 | 14.9 | 0.8 | 8.2 |
| MindSpeed | MindSpeed-Bridge | master | 12.8 | 8.3 | 8.3 | - | - | - | - | 13.3 | 8.1 | 8.1 | 0.0 | 2.4 |
| MindSpeed | slime-ascend | main | 25.1 | 12.3 | 10.0 | - | - | - | - | 24.4 | 21.4 | 21.4 | 0.0 | 2.1 |
| MindSpeed | modelzoo-GPL | master | 2.3 | 2.1 | 2.1 | - | - | - | - | - | - | - | - | 2.3 |
| MindSpeed | MindSpeed-RL | master | 0.4 | 0.4 | 0.4 | - | - | - | - | - | - | - | - | - |
| MindSpeed | apex | master | 4.4 | 4.4 | 4.4 | 1.8 | 1.8 | 1.8 | 0.0 | - | - | - | - | 2.5 |
| MindSpore | mindformers | master | 30.7 | 18.0 | 17.2 | - | - | - | - | 52.7 | 31.8 | 33.5 | 7.3 | 2.6 |
| MindStudio | msmodeling | master | 24.8 | 10.7 | 6.5 | 1.6 | 1.9 | 1.5 | 0.1 | 28.0 | 4.1 | 11.4 | 2.0 | 2.5 |
| MindStudio | msinsight | master | 35.1 | 27.7 | 33.3 | 12.2 | 13.7 | 12.8 | 1.3 | 18.1 | 16.2 | 16.5 | 0.2 | 2.6 |
| MindStudio | msmodelslim | master | 30.1 | 16.6 | 19.9 | 1.5 | 2.6 | 1.8 | 0.0 | 22.6 | 19.0 | 19.1 | 5.7 | 2.6 |
| MindStudio | msprobe | master | 12.9 | 7.9 | 8.2 | 0.8 | 4.0 | 2.3 | 0.1 | 7.3 | 5.9 | 6.1 | 0.4 | 2.4 |
| MindStudio | msprof | master | 17.7 | 7.8 | 1.9 | 3.4 | 5.8 | 4.4 | 0.0 | 10.8 | 9.2 | 9.6 | 1.1 | 2.5 |
| MindStudio | mssanitizer | master | 33.3 | 18.7 | 19.8 | 13.8 | 19.1 | 14.0 | 0.0 | 4.9 | 3.9 | 4.9 | 12.0 | 2.5 |
| MindStudio | msopcom | master | 17.7 | 13.0 | 11.8 | 2.3 | 2.4 | 2.0 | 4.8 | 5.8 | 5.7 | 5.5 | 0.0 | 2.5 |
| MindStudio | mspti | master | 12.4 | 6.7 | 6.3 | 2.1 | 2.9 | 2.2 | 0.1 | 7.0 | 6.2 | 6.2 | 1.0 | 2.4 |
| MindStudio | msagent | master | 4.7 | 3.6 | 4.5 | 0.8 | 1.0 | 0.8 | 0.0 | - | - | - | - | 2.4 |
| MindStudio | msmemscope | master | 32.1 | 19.6 | 17.8 | 6.9 | 18.3 | 10.4 | 0.0 | 16.4 | 14.3 | 14.8 | 9.2 | 2.6 |
| MindStudio | msmonitor | master | 9.6 | 4.9 | 2.0 | 3.2 | 7.0 | 4.1 | 0.1 | 9.9 | 9.2 | 9.4 | 1.1 | 2.2 |
| MindStudio | msopprof | master | 24.2 | 13.3 | 14.3 | 9.4 | 10.2 | 9.0 | 0.1 | 10.7 | 9.5 | 9.7 | 4.9 | 2.5 |
| MindStudio | msserviceprofiler | master | 10.8 | 5.3 | 5.0 | 3.1 | 11.1 | 5.7 | 0.0 | 6.4 | 2.7 | 3.2 | 0.0 | 2.6 |
| MindStudio | msdebug | master | 41.3 | 31.8 | 39.3 | 8.2 | 16.9 | 11.4 | 1.4 | 18.0 | 9.9 | 11.5 | 0.4 | 9.7 |
| MindStudio | msboost | master | 3.3 | 2.6 | 3.0 | - | - | - | - | - | - | - | - | 2.6 |
| MindStudio | msprof-analyze | master | 12.4 | 7.2 | 8.0 | 1.6 | 1.9 | 1.7 | 0.0 | 10.3 | 8.9 | 9.0 | 5.6 | 2.5 |
| MindStudio | mstt | master | 2.9 | 2.1 | 2.3 | - | - | - | - | 2.8 | 2.2 | 2.0 | 0.1 | 2.4 |
| MindStudio | msit | master | 4.6 | 3.0 | 2.8 | - | - | - | - | 1.9 | 0.3 | 1.0 | 0.0 | 2.3 |
| MindStudio | mskpp | master | 8.3 | 4.8 | 4.6 | 2.5 | 2.6 | 2.5 | 0.0 | 3.5 | 3.4 | 3.4 | 0.0 | 2.6 |
| MindStudio | mstx | master | 9.6 | 5.7 | 5.0 | 1.8 | 1.9 | 1.9 | 0.0 | 1.0 | 0.9 | 0.9 | 3.6 | 2.3 |
| MindStudio | mskl | master | 7.0 | 4.9 | 6.5 | 1.8 | 1.9 | 1.8 | 0.1 | 2.3 | 1.7 | 1.7 | 0.0 | 2.2 |
| MindStudio | msmodelslim_for_lingqu | master | 31.1 | 23.5 | 27.2 | 0.9 | 1.1 | 1.0 | 0.0 | 26.2 | 20.2 | 21.4 | 2.2 | 2.3 |
| MindStudio | msopgen | master | 7.1 | 6.5 | 6.7 | 1.8 | 2.1 | 1.9 | 0.0 | 2.4 | 1.9 | 1.8 | 0.0 | 2.3 |
| MindStudio | msot | master | 19.0 | 8.8 | 1.5 | 19.4 | 19.4 | 19.4 | 0.0 | - | - | - | - | 2.4 |
| MindStudio | msoptuner | master | 10.7 | 9.9 | 9.9 | 4.5 | 4.5 | 4.5 | 0.1 | 2.5 | 2.3 | 2.3 | 0.0 | 2.3 |
| MindStudio | msprof_for_lingqu | master | 15.6 | 10.9 | 10.9 | 2.7 | 3.1 | 2.7 | 0.1 | 5.5 | 3.7 | 3.7 | 7.9 | 2.4 |
| MindStudio | mscommreport | master | 0.7 | 0.7 | 0.7 | - | - | - | - | - | - | - | - | - |
| MindStudio | msinsight_for_lingqu | master | 25.7 | 25.7 | 25.7 | 9.8 | 9.8 | 9.8 | 6.9 | 10.7 | 10.7 | 10.7 | 0.0 | 3.8 |
| MindStudio | msmonitor_for_lingqu | master | 5.0 | 5.0 | 5.0 | 2.1 | 2.1 | 2.1 | 0.0 | - | - | - | - | 2.2 |
| MindStudio | msopcom_for_lingqu | master | 3.7 | 3.7 | 3.7 | 0.7 | 0.7 | 0.7 | 0.0 | - | - | - | - | 2.3 |
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
| MindStudio | 33 |

---
*报告生成时间: 2026-06-30 15:36:14 | 数据来源: openlibing-ops repo-pr-pipeline (efficiencyAvgTime* 字段)*
