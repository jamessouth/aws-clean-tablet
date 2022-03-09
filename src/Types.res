type poolData
type poolDataInput = {
  @as("UserPoolId") userPoolId: string,
  @as("ClientId") clientId: string,
  @as("AdvancedSecurityDataCollectionFlag") advancedSecurityDataCollectionFlag: bool,
}

type attributeData
type attributeDataInput = {
  @as("Name") name: string,
  @as("Value") value: string,
}

type userData
type userDataInput = {
  @as("Username") username: string,
  @as("Pool") pool: poolData,
}




