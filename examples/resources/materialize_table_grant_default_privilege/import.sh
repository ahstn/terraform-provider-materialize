#Grants can be imported using the concatenation of GRANT DEFAULT, the grantee id of the role
#Optionally you can include the target id, database id and schema id. The privilege is required 
terraform import materialize_table_grant_default_privilege.example GRANT DEFAULT|TABLE|<grantee_id>|<target_role_id>|<database_id>|<schema_id>|<privilege>
