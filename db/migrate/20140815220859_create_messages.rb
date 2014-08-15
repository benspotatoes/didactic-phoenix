class CreateMessages < ActiveRecord::Migration
  def change
    create_table :messages do |t|
      t.string :token
      t.string :team_id
      t.string :channel_id
      t.string :channel_name
      t.string :timestamp
      t.string :user_id
      t.string :string
      t.string :user_name
      t.text :text
      t.string :trigger_word

      t.timestamps
    end
  end
end
