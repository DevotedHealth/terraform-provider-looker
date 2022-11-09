resource "looker_group_membership" "group_membership" {
  target_group_id = looker_group.group.id
  user_ids        = [looker_user.user1.id, looker_user.user2.id, looker_user.user3.id]
  group_ids       = [looker_group.group1.id, looker_group.group2.id]
}
