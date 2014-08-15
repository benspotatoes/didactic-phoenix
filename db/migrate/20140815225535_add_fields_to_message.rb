class AddFieldsToMessage < ActiveRecord::Migration
  def change
    add_column :messages, :team_domain, :string
    add_column :messages, :service_id, :string
  end
end
