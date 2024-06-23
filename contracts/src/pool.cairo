use starknet::ContractAddress;

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
mod Pool {
    use starknet::ContractAddress;
    use openzeppelin::access::ownable::OwnableComponent;
    use super::IPool;
    use openzeppelin::token::erc20::dual20::DualCaseERC20;
    use openzeppelin::token::erc20::interface::IERC20;
    use openzeppelin::token::erc20::interface::IERC20DispatcherTrait;
    use openzeppelin::token::erc20::interface::IERC20Dispatcher;
    use starknet::get_caller_address;
    use starknet::get_contract_address;

    component!(path: OwnableComponent, storage: ownable, event: OwnableEvent);

    #[abi(embed_v0)]
    impl OwnableMixinImpl = OwnableComponent::OwnableMixinImpl<ContractState>;
    impl OwnableInternalImpl = OwnableComponent::InternalImpl<ContractState>;

    #[storage]
    struct Storage {
        #[substorage(v0)]
        ownable: OwnableComponent::Storage,
        token0: ContractAddress,// base token - ETH
        token1: ContractAddress,// Launched token
        reserve0: u256,
        reserve1: u256,
        totalSupply: u256
    }

    #[event]
    #[derive(Drop, starknet::Event)]
    enum Event {
        #[flat]
        OwnableEvent: OwnableComponent::Event,
        AddLiquidityEvent: AddLiquidity,
        SwapEvent: Swap,
    }

    #[derive(Drop, starknet::Event)]
    struct AddLiquidity {
        #[key]
        reserve1: u256,
        reserve0: u256,
        // token0: ContractAddress,// base token - ETH
        // token1: ContractAddress,// Launched token
        // totalSupply: u256
    }

    #[derive(Drop, starknet::Event)]
    struct Swap {
        #[key]
        tokenIn: ContractAddress,
        tokenOut: ContractAddress,
        amountIn: u256,
        amountOut: u256,
        reserve1: u256,
        reserve0: u256,
    }

    //TODO: Base token(ETH) can be part of the constructor as well
    #[constructor]
    fn constructor(ref self: ContractState, owner: ContractAddress, token0: ContractAddress) {
        self.ownable.initializer(owner);
        self.token0.write(token0);
        self.token1.write(token0);
        self.reserve0.write(0_u256);
        self.reserve1.write(0_u256);
        self.totalSupply.write(0_u256);
    }

    const ETH_ADDRESS: felt252 = 0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7;
    
    
    #[abi(embed_v0)]
    impl PoolImpl of IPool<ContractState> {
        //TODO: add amount0, amount1 : u256, hardcoding for now since testing with u256 is difficult
        fn add_liquidity(ref self: ContractState, amount0: u256, amount1: u256){
            self.ownable.assert_only_owner();
            let src = get_caller_address();
            let dst = get_contract_address();
            IERC20Dispatcher { contract_address: self.token1.read().try_into().unwrap() }.transfer_from(src, dst, amount1);

            let shares_to_allocate = if self.totalSupply.read() == 0_u256{
                sqrt(amount0 * amount1)
            }else{
                let shares_one = (amount0*self.totalSupply.read())/self.reserve0.read();
                let shares_two = (amount1*self.totalSupply.read())/self.reserve1.read();

                if shares_one > shares_two{
                    shares_two
                }else{
                    shares_one
                }
            };

            assert!(shares_to_allocate > 0, "invalid share allocation");
            let new_reserves_0 = self.reserve0.read() + amount0;
            let new_reserves_1 = self.reserve1.read() + amount1;

            self.reserve0.write(new_reserves_0);
            self.reserve1.write(new_reserves_1);
            //Grant shares now to the liquidity provider?
            self.emit(AddLiquidity{reserve0: new_reserves_0, reserve1: new_reserves_1});
        }

        fn swap(ref self: ContractState, amountIn: u256, style: u8){
            if style == 0 {
                self._buy(amountIn, true);
            }else{
                self._sell(amountIn, true);
            }
        }
        
        fn quote(self: @ContractState, amountIn: u256, style: u8) -> u256{
            if style == 0 {
                self._quote(self.reserve1.read(), self.reserve0.read(), amountIn)
            }else{
                self._quote(self.reserve0.read(), self.reserve1.read(), amountIn)
            }
        }

        fn get_balance(self: @ContractState) -> (u256, u256) {
            let currentAddress = get_contract_address();
            let tokenBal = IERC20Dispatcher { contract_address: self.token0.read().try_into().unwrap() }.balance_of(currentAddress);
            let ethBal = IERC20Dispatcher { contract_address: ETH_ADDRESS.try_into().unwrap() }.balance_of(currentAddress);
            (ethBal, tokenBal)
        }

        fn get_reserves(self: @ContractState) -> (u256, u256){
            (self.reserve0.read(), self.reserve1.read())
        }

        fn get_pool_token(self: @ContractState) -> ContractAddress{
            self.token1.read()
        }
        
        fn add_gas(ref self: ContractState){
            // let contractAddress = get_contract_address();
        }

        fn get_token_approvals(self: @ContractState) -> (u256, u256){
            let contractAddress = get_contract_address();
            let callerAddress = get_caller_address();
            let a1 = IERC20Dispatcher { contract_address: ETH_ADDRESS.try_into().unwrap() }.allowance(callerAddress, contractAddress);
            let a2 = IERC20Dispatcher { contract_address: self.token1.read().try_into().unwrap() }.allowance(callerAddress, contractAddress);
            (a1, a2)
        }
    }

    #[generate_trait]
    impl InternalFunctions of InternalFunctionsTrait {
        fn _buy(
            ref self: ContractState,
            ethAmount: u256, 
            toSwap: bool, 
        ) -> u256 {
            let callerAddress = get_caller_address();
            let currentAddress = get_contract_address();
            let reserveIn = self.reserve0.read();
            let reserveOut = self.reserve1.read();

            let amountOut: u256 = self._quote(reserveOut, reserveIn, ethAmount);
            if toSwap{
                IERC20Dispatcher { contract_address: ETH_ADDRESS.try_into().unwrap()}.transfer_from(callerAddress, currentAddress, ethAmount);
                IERC20Dispatcher { contract_address: self.token0.read()}.transfer(callerAddress, amountOut);
                self.reserve0.write(self.reserve0.read() + ethAmount);
                self.reserve1.write(self.reserve1.read() - amountOut);
                let newReserve0 = IERC20Dispatcher { contract_address: ETH_ADDRESS.try_into().unwrap()}.balance_of(currentAddress);
                let newReserve1 = IERC20Dispatcher { contract_address: self.token0.read()}.balance_of(currentAddress);
                self._update(newReserve0, newReserve1);
                self.emit(Swap{
                    tokenIn: ETH_ADDRESS.try_into().unwrap(),
                    tokenOut: self.token0.read(),
                    amountIn: ethAmount,
                    amountOut: amountOut,
                    reserve0: self.reserve0.read(),
                    reserve1: self.reserve1.read(),
                })
            }
            amountOut
        }

        fn _sell(
            ref self: ContractState,
            tokenAmount: u256, 
            toSwap: bool, 
        ) -> u256 {
            let callerAddress = get_caller_address();
            let currentAddress = get_contract_address();
            let reserveIn = self.reserve1.read();
            let reserveOut = self.reserve0.read();

            let amountOut: u256 = self._quote(reserveOut, reserveIn, tokenAmount);
            if toSwap{
                IERC20Dispatcher { contract_address: self.token0.read()}.transfer_from(callerAddress, currentAddress, tokenAmount);
                IERC20Dispatcher { contract_address: ETH_ADDRESS.try_into().unwrap()}.transfer(callerAddress, amountOut);
                self.reserve0.write(self.reserve0.read() - amountOut);
                self.reserve1.write(self.reserve1.read() + tokenAmount);
                
                let newReserve0 = IERC20Dispatcher { contract_address: ETH_ADDRESS.try_into().unwrap()}.balance_of(currentAddress);
                let newReserve1 = IERC20Dispatcher { contract_address: self.token0.read()}.balance_of(currentAddress);
                self._update(newReserve0, newReserve1);
                self.emit(Swap{
                    tokenOut: ETH_ADDRESS.try_into().unwrap(),
                    tokenIn: self.token0.read(),
                    amountIn: tokenAmount,
                    amountOut: amountOut,
                    reserve0: self.reserve0.read(),
                    reserve1: self.reserve1.read(),
                })
            }
            amountOut
        }
 
        fn _quote(self: @ContractState, reserveOut: u256, reserveIn: u256, amountIn: u256) -> u256{
            (reserveOut * amountIn)/(reserveIn + amountIn)
        }

        fn _update(ref self: ContractState, newReserve0: u256, newReserve1: u256){
            self.reserve0.write(newReserve0);
            self.reserve1.write(newReserve1);
        }
    }

    fn sqrt(n: u256) -> u256{
        if n == 0 {
            return 0;
        }
        let mut guess = n / 2;
        let mut diff = guess;
        while diff > 1 {
            let new_guess = (guess + n / guess) / 2;
            if new_guess > guess {
                diff = new_guess - guess;
            } else {
                diff = guess - new_guess;
            }
            guess = new_guess;
        };
        return guess;
    }
}
