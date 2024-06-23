use starknet::ContractAddress;

#[starknet::interface]
trait IDeployer<TContractState> {
    fn deployERC20(ref self: TContractState, name: ByteArray, symbol: ByteArray, user: ContractAddress) -> ContractAddress;
}

#[starknet::interface]
trait IPool<TContractState> {
    fn add_liquidity(ref self: TContractState, amount0: u256, amount1: u256);
    fn quote(self: @TContractState, amountIn: u256, style: u8) -> u256;
    fn swap(ref self: TContractState, amountIn: u256, style: u8);
    fn get_balance(self: @TContractState) -> (u256, u256);
    fn get_reserves(self: @TContractState) -> (u256, u256);
    fn get_pool_token(self: @TContractState) -> ContractAddress;
    fn get_token_approvals(self: @TContractState) -> (u256, u256);
    fn add_gas(ref self: TContractState);
}

#[starknet::contract]
mod Deployer {
    use super::IDeployer;
    use super::IPool;
    use super::IPoolDispatcherTrait;
    use super::IPoolDispatcher;
    use openzeppelin::utils::interfaces::{IUniversalDeployerDispatcher, IUniversalDeployerDispatcherTrait};
    use starknet::ContractAddress;
    use starknet::class_hash::class_hash_const;
    use openzeppelin::utils::serde::SerializedAppend;
    use openzeppelin::access::ownable::OwnableComponent;
    use starknet::get_caller_address;
    use starknet::get_contract_address;
    use openzeppelin::token::erc20::interface::IERC20;
    use openzeppelin::token::erc20::interface::IERC20DispatcherTrait;
    use openzeppelin::token::erc20::interface::IERC20Dispatcher;

    component!(path: OwnableComponent, storage: ownable, event: OwnableEvent);

    #[abi(embed_v0)]
    impl OwnableMixinImpl = OwnableComponent::OwnableMixinImpl<ContractState>;
    impl OwnableInternalImpl = OwnableComponent::InternalImpl<ContractState>;

    const UDC_ADDRESS: felt252 = 0x04a64cd09a853868621d94cae9952b106f2c36a3f81260f85de6696c6b050221;
    const ETH_ADDRESS: felt252 = 0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7;

    #[storage]
    struct Storage {
        #[substorage(v0)]
        ownable: OwnableComponent::Storage
    }

    #[event]
    #[derive(Drop, starknet::Event)]
    enum Event {
        #[flat]
        OwnableEvent: OwnableComponent::Event,
        TokenCreationEvent: TokenCreation,
        PoolCreationEvent: PoolCreation,
    }

    #[derive(Drop, starknet::Event)]
    struct TokenCreation {
        #[key]
        user: ContractAddress,
        token: ContractAddress,
    }

    #[derive(Drop, starknet::Event)]
    struct PoolCreation {
        #[key]
        user: ContractAddress,
        token: ContractAddress,
        pool: ContractAddress,
    }

    #[constructor]
    fn constructor(ref self: ContractState, owner: ContractAddress) {
        self.ownable.initializer(owner);
    }

    #[abi(embed_v0)]
    impl Deployer of IDeployer<ContractState> {
        // 1. Launch Token
        // 2. Launch Pool
        // 3. Give approval to the new token to be added by the pool
        // 4. Add some gas balance
        // 5. Call add liquidity function in the pool contract
        // 6. Give ETH approval as well? 
        fn deployERC20(ref self: ContractState, name: ByteArray, symbol: ByteArray, user: ContractAddress) -> ContractAddress {
            self.ownable.assert_only_owner();

            let dispatcher = IUniversalDeployerDispatcher {
                contract_address: UDC_ADDRESS.try_into().unwrap()
            };

            /////////////////////////// STEP 1 /////////////////////////////////////////////////////////////
            let class_hash = class_hash_const::<0x040b9e69e14ddc34a98ec8133c80807c144b818bc6cbf5a119d8f62535258142>();
            let salt: felt252 = 'SAFTC8989';
            let from_zero = false;

            let totalSupply: u256 = 1000000000000000000000000000;
            let mut calldata = array![];
            calldata.append_serde(name);
            calldata.append_serde(symbol);
            calldata.append_serde(totalSupply);
            calldata.append_serde(get_contract_address());
            calldata.append_serde(self.ownable.owner());

            let tokenAddress = dispatcher.deploy_contract(class_hash, salt, from_zero, calldata.span());
            self.emit(TokenCreation{user: user, token: tokenAddress});
            /////////////////////////// STEP 2 /////////////////////////////////////////////////////////////
            let pool_class_hash = class_hash_const::<0x0621dc859cf188aa6b6cd0826fba497423d744bc4442fc5758acf8128ccfea07>();
            let pool_salt: felt252 = 'SAFTC8980';

            let mut pool_calldata = array![];
            pool_calldata.append_serde(get_contract_address());
            pool_calldata.append_serde(tokenAddress);
            let poolAddress = dispatcher.deploy_contract(pool_class_hash, pool_salt, from_zero, pool_calldata.span());
            self.emit(PoolCreation{user: user, token: tokenAddress, pool: poolAddress});
            /////////////////////////// STEP 3 /////////////////////////////////////////////////////////////
            IERC20Dispatcher { contract_address: tokenAddress }.approve(poolAddress, totalSupply);
            /////////////////////////// STEP 4 /////////////////////////////////////////////////////////////
            IERC20Dispatcher { contract_address: ETH_ADDRESS.try_into().unwrap()}.transfer_from(
                self.ownable.owner(), 
                poolAddress, 
                2500000000000000_u256
            );
            /////////////////////////// STEP 5 /////////////////////////////////////////////////////////////
            IPoolDispatcher {contract_address: poolAddress}.add_liquidity(1500000000000000000_u256, 1000000000000000000000000000_u256);
            tokenAddress
        }
    }
}
