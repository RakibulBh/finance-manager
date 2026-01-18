"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Modal } from "@/components/ui/modal";
import { useAccounts } from "@/hooks/useAccounts";
import { useCreateTransaction } from "@/hooks/useTransactions";
import React, { useState } from "react";

interface AddTransactionModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export const AddTransactionModal = ({ isOpen, onClose, onSuccess }: AddTransactionModalProps) => {
  const [formData, setFormData] = useState({
    account_id: "",
    name: "",
    amount: "",
    date: new Date().toISOString().split("T")[0],
    merchant_name: "",
  });

  const { data: accounts = [] } = useAccounts();
  const { mutate: createTransaction, isPending: loading, error: mutationError } = useCreateTransaction();

  // Derived state for the effective account ID to avoid setState in effect
  const effectiveAccountId = formData.account_id || accounts[0]?.id || "";

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    createTransaction(
      {
        ...formData,
        account_id: effectiveAccountId,
        amount: parseFloat(formData.amount),
      },
      {
        onSuccess: () => {
          onSuccess();
          onClose();
          // Reset form
          setFormData({
            account_id: accounts[0]?.id || "",
            name: "",
            amount: "",
            date: new Date().toISOString().split("T")[0],
            merchant_name: "",
          });
        },
      }
    );
  };

  if (!isOpen) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Add Transaction">
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-1.5">
          <label className="block text-sm font-medium text-zinc-400 ml-1" htmlFor="account-select">Account</label>
          <select
            id="account-select"
            title="Select Account"
            className="w-full bg-brand-muted border border-zinc-800 rounded-2xl px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-brand-lime/20 transition-all appearance-none"
            value={effectiveAccountId}
            onChange={(e) => setFormData({ ...formData, account_id: e.target.value })}
            required
          >
            {accounts.map((acc) => (
              <option key={acc.id} value={acc.id}>
                {acc.name} (${acc.balance.toLocaleString()})
              </option>
            ))}
          </select>
        </div>

        <Input
          label="Label / Description"
          placeholder="e.g. Starbucks Coffee"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          required
        />

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="Amount ($)"
            type="number"
            step="0.01"
            placeholder="0.00"
            value={formData.amount}
            onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
            required
          />
          <Input
            label="Date"
            type="date"
            value={formData.date}
            onChange={(e) => setFormData({ ...formData, date: e.target.value })}
            required
          />
        </div>

        <Input
          label="Merchant Name (Optional)"
          placeholder="e.g. Starbucks"
          value={formData.merchant_name}
          onChange={(e) => setFormData({ ...formData, merchant_name: e.target.value })}
        />

        {mutationError && (
          <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-red-500 text-sm">
            {(mutationError as Error).message}
          </div>
        )}

        <Button type="submit" className="w-full py-4" disabled={loading}>
          {loading ? "Recording..." : "Record Transaction"}
        </Button>
      </form>
    </Modal>
  );
};
