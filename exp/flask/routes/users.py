from flask import Blueprint, request, jsonify

users_bp = Blueprint('users', __name__)

# GET all users
@users_bp.route('/')
def get_users():
    page = request.args.get('page')
    limit = request.args.get('limit')
    return jsonify({'users': []})

# GET user by ID
@users_bp.route('/<int:user_id>')
def get_user(user_id):
    return jsonify({'id': user_id})

# POST create user
@users_bp.route('/', methods=['POST'])
def create_user():
    name = request.form.get('name')
    email = request.form.get('email')
    return jsonify({'id': 1})

# PUT update user
@users_bp.route('/<int:user_id>', methods=['PUT'])
def update_user(user_id):
    name = request.form.get('name')
    return jsonify({'id': user_id})

# DELETE user
@users_bp.route('/<int:user_id>', methods=['DELETE'])
def delete_user(user_id):
    return jsonify({'message': 'User deleted'})

# GET user's posts
@users_bp.route('/<int:user_id>/posts')
def get_user_posts(user_id):
    return jsonify({'posts': []})

