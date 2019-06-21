# Differences between Azure China Cloud and Azure Public Cloud

Below are known differences between Azure China Cloud and Azure Public Cloud(there may exist unknown differences not listed in this page):

- AppInsights

  Not supported in Azure China Cloud.

- MSSQL

  - For sql server created in chinanorth and chinaeast, vCore-based databases created upon it do not support Gen5 hardware.(OSBA will auto-switch to Gen4)

- MySQL

  - Only chinanorth and chinaeast is available. DO NOT use chinanorth2/chinaeast2.
  - Memory optimized plan is not supported in Azure China Cloud.
  - Gen5 hardware is not supported. (OSBA will auto-switch to Gen4)

- PostgreSQL

  - Memory optimized plan is not supported in chinanorth and chinaeast.
  - Gen5 hardware is not supported in chinanorth and chinaeast. (OBSA will auto-switch to Gen4)

- Storage

  ZRS account type is not supported in Azure China Cloud.

- Text Analytics

  Not supported in Azure China Cloud.