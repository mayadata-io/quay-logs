**1.** **How to generate quay auth token?**

In order to be able to download quay logs and quay repositories popularity related data, you will need to have quay auth token. Here are steps to generate one: (for internal application use)

1. Login to your [quay.io](https://quay.io) account and go to your organization.

2. Click on `Applications` on left sidebar.

3. Create an oauth application by clicking `New Application`.

4. Now click on the application. A new screen will come up with 4 tabs in left sidebar.

5. Now click on `Generate Token` tab and select `Administer Repositories` scope & click `Generate Access Token`.

A new access token will be generate on behalf of your account id. Access tokens for the quay.io are long-lived and do not expire.